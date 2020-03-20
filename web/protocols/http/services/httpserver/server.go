package httpserver

import (
	"context"
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/published"
	"gerrit-share.lan/go/actors/starter"
	"gerrit-share.lan/go/basicerrors"
	"gerrit-share.lan/go/debug"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectable"
	"gerrit-share.lan/go/inspect/inspectables"
	"gerrit-share.lan/go/inspect/inspectwrappers"
	"gerrit-share.lan/go/inspect/json/fromjson"
	"gerrit-share.lan/go/inspect/json/tojson"
	"gerrit-share.lan/go/interfaces"
	"gerrit-share.lan/go/utils/flags"
	"gerrit-share.lan/go/utils/maps"
	"gerrit-share.lan/go/web/auth"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/actorutils"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/cookies"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/frombytes"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/metadata"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/replies"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/tobytes"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/jsonrpc"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func NewHttpServer(hostPort flags.HostPort, name string, dir string, config *Config) *httpServer {
	server := new(httpServer)
	server.name = name
	server.dir = dir
	server.config = config
	server.server = http.Server{
		Addr:    hostPort.String(),
		Handler: server}
	return server
}

type httpServer struct {
	actors.Actor
	endpoints   HttpRestEndpoints
	server      http.Server
	serverActor actors.ActorService
	name        string
	dir         string
	authHandle  auth.Handle
	config      *Config
	noAuth      bool
}

func (h *httpServer) MakeBehaviour() actors.Behaviour {
	log.Println(h.name, "started")
	var starterHandle starter.Handle
	starterHandle.Acquire(h, starterHandle.DependOn, h.Quit)
	h.authHandle.Acquire(h, nil, func(err error) {
		if !errors.Is(err, actors.ErrNotGonnaHappen) {
			h.Quit(err)
		}
		h.noAuth = true
	})
	h.loadEndpoints()
	h.loadConfig()
	h.subscribeForActors()
	h.registerOwnEndpoints()
	var behaviour actors.Behaviour
	behaviour.Name = DefaultHttpServerName
	behaviour.AddCommand(new(urlPath), func(cmd interface{}) (actors.Response, error) {
		return h.generateDoc(cmd.(*urlPath))
	}).ResultString()
	behaviour.AddCommand(new(httpRequest), func(cmd interface{}) (actors.Response, error) {
		return h.serveRequest(cmd.(*httpRequest))
	}).Result(new(replies.HttpResponse))
	behaviour.AddCommand(new(editEndpoint), func(cmd interface{}) (actors.Response, error) {
		return nil, h.endpoints.Edit(cmd.(*editEndpoint))
	})
	behaviour.AddCommand(new(getEndpoints), func(cmd interface{}) (actors.Response, error) {
		return h.getEndpoints(cmd.(*getEndpoints)), nil
	}).Result(new(endpointInfoByName))
	behaviour.AddCommand(new(editEndpointByName), func(cmd interface{}) (actors.Response, error) {
		return nil, h.endpoints.EditByName(cmd.(*editEndpointByName))
	})
	h.SetPanicProcessor(h.onPanic)
	h.SetFinishedServiceProcessor(h.onServiceFinished)
	h.SetExitProcessor(h.onExit)
	h.serverActor = h.System().RunAsyncSimple(func() error {
		log.Println("listen and serve started")
		h.server.ListenAndServe()
		log.Println("listen and serve shutdown")
		return nil
	})
	h.Monitor(h.serverActor)
	return behaviour
}

func (h *httpServer) registerOwnEndpoints() {
	h.processEndpointMessage(urlPathName, new(urlPath), inspectwrappers.NewStringValue(""), h.Service(), ".doc")
	h.processEndpointMessage(editEndpointName, new(editEndpoint), nil, h.Service(), "edit.api")
	h.processEndpointMessage(getEndpointsName, new(getEndpoints), new(endpointInfoByName), h.Service(), "get.api")
	h.processEndpointMessage(editEndpointByNameName, new(editEndpointByName), nil, h.Service(), "editbyname.api")
}

func (h *httpServer) subscribeForActors() {
	input := actors.NewSimpleCallbackStreamInput(func(data inspect.Inspectable) error {
		array := data.(*actors.ActorsArray)
		for _, a := range *array {
			actor := a
			h.Monitor(actor)
			h.SendRequest(actor, actors.GetInfo{},
				actors.OnReply(func(reply interface{}) {
					endpoints := reply.(*actors.ActorCommands)
					for _, command := range endpoints.Commands {
						h.processEndpointMessage(command.Command.TypeName(), command.Command.Sample(),
							command.Result.Sample(), actor, prepareEndpointName(command.Command.TypeName(), endpoints.Name))
					}
				}))
		}
		return nil
	}, func(base *actors.StreamInputBase) {
		base.RequestData(new(actors.ActorsArray), 10)
	})
	published.Subscribe(h, input, h.Quit)
}

func (h *httpServer) Shutdown() error {
	h.onExit()
	log.Println(h.name, "shut down")
	return nil
}

func (h *httpServer) onPanic(err errors.StackTraceError) {
	log.Println("panic:", err, err.StackTrace())
	h.onExit()
	h.Quit(err)
}

func (h *httpServer) onServiceFinished(service actors.ActorService, err error) {
	if service == h.serverActor {
		log.Println("http server finished with error:", err)
		h.Quit(err)
	}
	debug.Printf("service finished: %p\n", service)
	h.endpoints.RemoveByDestination(service)
}

func (h *httpServer) onExit() {
	err := h.saveEndpoints()
	if err != nil {
		log.Println("failed to save endpoints:", err)
	}
	err = h.server.Shutdown(context.Background())
	if err != nil {
		log.Println("error while shutdown:", err)
	}
}

func (h *httpServer) getEndpoints(command *getEndpoints) *endpointInfoByName {
	result := &endpointInfoByName{}
	for name, info := range h.endpoints.info {
		if command.edited == info.changed {
			newInfo := &endpointInfo{changed: true}
			if command.edited {
				newInfo.update(info.Resource, info.Method, info.HttpMethod)
			} else {
				parsed := ParseEndpoint(name)
				newInfo.update(parsed.Name, parsed.Method, defaultHttpMethod)
			}
			(*result)[name] = newInfo
		}
	}
	return result
}

func (h *httpServer) loadEndpoints() {
	content, err := ioutil.ReadFile(path.Join(h.dir, endpointsFileName))
	if err != nil {
		log.Println("failed to load endpoints:", err)
	}
	reader := inspect.NewGenericInspector(fromjson.NewInspector(content, 0))
	h.endpoints.Inspect(reader)
}

func (h *httpServer) loadConfig() {
	// TODO: implement real config storage
	if h.config == nil {
		h.config = &defaultConfig
	}
}

func (h *httpServer) saveEndpoints() error {
	file, err := os.Create(path.Join(h.dir, endpointsFileName))
	if err != nil {
		return err
	}
	defer file.Close()
	inspector := &tojson.Inspector{}
	serializer := inspect.NewGenericInspector(inspector)
	h.endpoints.info.Inspect(serializer)
	if serializer.GetError() != nil {
		return serializer.GetError()
	}
	_, err = file.Write(inspector.Output())
	return nil
}

func (h *httpServer) processEndpointMessage(commandName string, commandSample inspect.Inspectable,
	resultSample inspect.Inspectable, service actors.ActorService, name string) {

	commandMetaData, err := metadata.MakeCommandMetaData(commandSample)
	if err != nil {
		log.Printf("failed to parse command sample: %v\nendpoint has been ignored by http-server\n", err)
		return
	}
	commandMetaData.Description = inspectables.GetDescription(commandName)

	if metadata.IsNested(commandMetaData) {
		log.Printf("nested command has been detected during registration of endpoint: %v\nendpoint has been ignored by http-server\n", commandName)
		return
	}

	result, err := MakeResultInfo(resultSample)
	if err != nil {
		log.Printf("failed to parse result sample: %v\nendpoint has been ignored by http-server\n", err)
		return
	}

	h.endpoints.Add(name, service, inspectables.Get(commandName), result, commandMetaData)
}

func (h *httpServer) getEndpointInfo(request *http.Request) ([]byte, endpointInfo, error) {
	var endpoint endpointInfo
	path, method, methods := h.endpoints.FindEndpointMethods(request.URL.Path, '/')
	debug.Println("path:", request.URL.Path)
	debug.Println(string(path), string(method))
	if methods == nil {
		return path, endpoint, ErrResourceNotFound
	}
	httpMethods, ok := methods[string(method)]
	if !ok {
		httpMethods, ok = methods[""]
		if !ok {
			return path, endpoint, ErrMethodNotAllowed
		}
		path = method
	}
	info, ok := httpMethods[strings.ToLower(request.Method)]
	if !ok || info.Dest == nil {
		return path, endpoint, ErrHttpMethodNotAllowed
	}
	endpoint = *info
	if endpoint.CommandMetaData.PathIndex == metadata.NoPath && len(path) > 0 {
		return path, endpoint, ErrPathNotAllowed
	}
	return path, endpoint, nil
}

func logFailedToSerializeErr(endpoint string, err error) {
	log.Printf("failed to serialize response from %v. %v", endpoint, err)
}

func setConnectionId(command inspect.Inspectable, connectionId uuid.UUID) {
	authCommand, ok := command.(auth.Info)
	if !ok {
		return
	}
	authCommand.SetConnectionId(auth.Id(connectionId))
}

func (h *httpServer) processSingleCommand(command inspect.Inspectable, endpoint endpointInfo,
	sessionInfo cookies.CookieInfo) actors.Response {
	promise := &replies.HttpResponsePromise{}
	setConnectionId(command, sessionInfo.SessionId)
	canceller := h.SendRequest(endpoint.Dest, command,
		actors.OnReply(func(reply interface{}) {
			result := inspectable.NewGenericValue(endpoint.ResultInfo.TypeId)
			result.SetValue(reply)
			bytes, err := tobytes.ToBytes(result)
			if err != nil {
				logFailedToSerializeErr(endpoint.OriginalName, err)
				promise.Fail(err)
				return
			}
			response := replies.NewHttpResponse(sessionInfo)
			response.Response = bytes
			promise.Deliver(response)
		}).OnError(promise.Fail))
	promise.OnCancel(canceller.Cancel)
	return promise
}

func (h *httpServer) processCommandBatch(commands []inspect.Inspectable, endpoint endpointInfo,
	sessionInfo cookies.CookieInfo) actors.Response {
	promise := &replies.HttpResponsePromise{}
	promiseCounter := len(commands)
	result := make(groupResponse, len(commands))
	cancellers := make([]interfaces.Canceller, len(commands))
	deliver := func() {
		promiseCounter--
		if promiseCounter == 0 {
			writer := &tojson.Inspector{}
			serializer := inspect.NewGenericInspector(writer)
			result.Inspect(serializer)
			if serializer.GetError() != nil {
				logFailedToSerializeErr(endpoint.OriginalName, serializer.GetError())
				promise.Fail(serializer.GetError())
				return
			}
			response := replies.NewHttpResponse(sessionInfo)
			response.Response = writer.Output()
			promise.Deliver(response)
		}
	}
	for i, command := range commands {
		index := i
		result[index].Result = inspectable.NewGenericValue(endpoint.ResultInfo.TypeId)
		setConnectionId(command, sessionInfo.SessionId)
		cancellers = append(cancellers, h.SendRequest(endpoint.Dest, command, actors.OnReply(func(reply interface{}) {
			result[index].Result.SetValue(reply)
			deliver()
		}).OnError(func(err error) {
			result[index].Err = err
			deliver()
		})))
	}
	promise.OnCancel(func() {
		for _, canceller := range cancellers {
			canceller.Cancel()
		}
	})
	return promise
}

func (h *httpServer) processJsonRequest(request *http.Request, endpoint endpointInfo,
	path []byte, sessionInfo cookies.CookieInfo) (actors.Response, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	var acceptValues bool
	params := endpoint.CommandMetaData.UnderlyingValues
	if len(params) == 1 || len(params) == 2 && !(endpoint.CommandMetaData.PathIndex == metadata.NoPath) {
		acceptValues = true
	}

	inspector := frombytes.NewProxyInspector(path, body, 0, acceptValues)
	deserializer := inspect.NewGenericInspector(inspector)
	objectCommands := &commandBatch{generator: endpoint.CommandGenerator}
	objectCommands.Inspect(deserializer)
	if deserializer.GetError() != nil {
		return nil, deserializer.GetError()
	}
	if inspector.IsBatch {
		return h.processCommandBatch(objectCommands.commands, endpoint, sessionInfo), nil
	}
	return h.processSingleCommand(objectCommands.commands[0], endpoint, sessionInfo), nil
}

func (h *httpServer) processUnknownRequest(request *http.Request, endpoint endpointInfo,
	path []byte, sessionInfo cookies.CookieInfo) (actors.Response, error) {
	bodyPart := make([]byte, 1)
	n, _ := request.Body.Read(bodyPart)
	if n > 0 {
		return nil, ErrNoContentType
	}
	command, err := commandFromPath(endpoint, path)
	if err != nil {
		return nil, err
	}
	return h.processSingleCommand(command, endpoint, sessionInfo), nil
}

func (h *httpServer) processRequest(request *http.Request, endpoint endpointInfo,
	path []byte, sessionInfo cookies.CookieInfo) (actors.Response, error) {
	if request.Method == http.MethodGet {
		command, err := commandFromQueryString(request, endpoint, path)
		if err != nil {
			return nil, err
		}
		return h.processSingleCommand(command, endpoint, sessionInfo), nil
	}
	switch contentType := request.Header.Get("Content-Type"); contentType {
	case "application/x-www-form-urlencoded":
		command, err := commandFromQueryString(request, endpoint, path)
		if err != nil {
			return nil, err
		}
		return h.processSingleCommand(command, endpoint, sessionInfo), nil
	case "application/json":
		return h.processJsonRequest(request, endpoint, path, sessionInfo)
	default:
		if strings.HasPrefix(contentType, "multipart/form-data") {
			command, err := commandFromFormData(request, endpoint, path)
			if err != nil {
				return nil, err
			}
			return h.processSingleCommand(command, endpoint, sessionInfo), nil
		}
		if len(contentType) > 0 {
			return nil, ErrUnsupportedContentType
		}
		return h.processUnknownRequest(request, endpoint, path, sessionInfo)
	}
}

func (h *httpServer) processRpcCommand(command requestWithDestination, sessionInfo cookies.CookieInfo) (actors.Response, error) {
	rpcResponse := replies.NewHttpResponse(sessionInfo)
	if command.err != nil {
		response := jsonrpc.NewResponse(command.responseType)
		response.Result.Err = command.err
		bytes, err := tobytes.ToBytes(response)
		if err != nil {
			return nil, err
		}
		rpcResponse.Response = bytes
		return rpcResponse, nil
	}
	promise := &replies.HttpResponsePromise{}
	response := jsonrpc.NewResponse(command.responseType)
	writer := &tojson.Inspector{}
	serializer := inspect.NewGenericInspector(writer)
	deliver := func() {
		response.Inspect(serializer)
		if serializer.GetError() != nil {
			logFailedToSerializeErr(command.request.Method, serializer.GetError())
			promise.Fail(jsonrpc.Describe(serializer.GetError(), jsonrpc.ErrInternalError))
			return
		}
		rpcResponse.Response = writer.Output()
		promise.Deliver(rpcResponse)
	}
	setConnectionId(command.request.Params, sessionInfo.SessionId)
	canceller := h.SendRequest(command.destination, command.request.Params, actors.OnReply(func(reply interface{}) {
		response.Result.Result.SetValue(reply)
		deliver()
	}).OnError(func(err error) {
		response.Result.Err = jsonrpc.Describe(err, jsonrpc.ErrInvalidParams)
		deliver()
	}))
	promise.OnCancel(canceller.Cancel)
	return promise, nil
}

func (h *httpServer) processRpcCommandBatch(commands []requestWithDestination, sessionInfo cookies.CookieInfo) (actors.Response, error) {
	promise := &replies.HttpResponsePromise{}
	rpcResponse := replies.NewHttpResponse(sessionInfo)
	var validCommands int
	for _, command := range commands {
		if command.err == nil {
			validCommands++
		}
	}
	cancellers := make([]interfaces.Canceller, validCommands)
	result := make(jsonrpc.ResponseBatch, len(commands))
	deliver := func() {
		validCommands--
		if validCommands == 0 {
			output, err := tobytes.ToBytes(&result)
			if err != nil {
				promise.Fail(jsonrpc.Describe(err, jsonrpc.ErrInternalError))
				return
			}
			rpcResponse.Response = output
			promise.Deliver(rpcResponse)
		}
	}
	for i, command := range commands {
		index := i
		response := jsonrpc.NewResponse(command.responseType)
		if command.err != nil {
			response.Result.Err = command.err
			result[index] = response
			continue
		}
		setConnectionId(command.request.Params, sessionInfo.SessionId)
		cancellers = append(cancellers, h.SendRequest(command.destination, command.request.Params, actors.OnReply(func(reply interface{}) {
			response.Result.Result.SetValue(reply)
			result[index] = response
			deliver()
		}).OnError(func(err error) {
			response.Result.Err = jsonrpc.Describe(err, jsonrpc.ErrInvalidParams)
			result[index] = response
			deliver()
		})))
	}
	if validCommands == 0 {
		output, err := tobytes.ToBytes(&result)
		if err != nil {
			return nil, jsonrpc.Describe(err, jsonrpc.ErrInvalidParams)
		}
		rpcResponse.Response = output
		return rpcResponse, nil
	}
	promise.OnCancel(func() {
		for _, canceller := range cancellers {
			canceller.Cancel()
		}
	})
	return promise, nil
}

func (h *httpServer) processJsonRpcRequest(request *http.Request, sessionInfo cookies.CookieInfo) (actors.Response, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, jsonrpc.Describe(err, jsonrpc.ErrInternalError)
	}

	inspector := frombytes.NewProxyBatchInspector(body, 0)
	deserializer := inspect.NewGenericInspector(inspector)
	objectCommands := &rpcRequestBatch{endpoints: &h.endpoints}
	objectCommands.Inspect(deserializer)
	if deserializer.GetError() != nil {
		return nil, deserializer.GetError()
	}
	if inspector.IsBatch {
		return h.processRpcCommandBatch(objectCommands.data, sessionInfo)
	}
	return h.processRpcCommand(objectCommands.data[0], sessionInfo)
}

func (h *httpServer) serveRequest(cmd *httpRequest) (actors.Response, error) {
	sessionId, expiresAt, err := cookies.ParseCookie(&cmd.Request, h.config.JwtKey)
	var id uuid.UUID
	if err != nil && !errors.Is(err, cookies.ErrTokenExpired) {
		return nil, err
	}
	if err != nil || sessionId == nil {
		id = uuid.New()
		if !h.noAuth {
			h.authHandle.ConnectionEstablished(auth.Id(id))
		}
	} else {
		id, err = uuid.ParseBytes(sessionId)
		if err != nil {
			return nil, cookies.ErrIncorrectSessionToken
		}
	}
	sessionInfo := cookies.CookieInfo{
		SessionId:       id,
		ExpiresAt:       expiresAt,
		SessionReset:    h.config.SessionResetPercent,
		SessionDuration: h.config.SessionDuration,
	}
	if cmd.Header.Get("Content-Type") == "application/json-rpc" {
		return h.processJsonRpcRequest(&cmd.Request, sessionInfo)
	}
	path, endpoint, err := h.getEndpointInfo(&cmd.Request)
	if err != nil {
		return nil, err
	}
	return h.processRequest(&cmd.Request, endpoint, path, sessionInfo)
}

func (h *httpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	var cancelled bool
	h.System().Become(actorutils.NewShutdownableActor(request.Context().Done(),
		func() {
			writer.WriteHeader(http.StatusBadGateway)
			cancelled = true
		},
		func(actor *actors.Actor) actors.Behaviour {
			var b actors.Behaviour
			actor.SendRequest(h.Service(), &httpRequest{Request: *request},
				actors.OnReply(func(reply interface{}) {
					replyTyped := reply.(*replies.HttpResponse)
					cookies.AddCookie(writer, replyTyped.SessionId, h.config.JwtKey, h.config.SessionDuration)
					writer.Write(replyTyped.Response)
				}).OnError(func(err error) {
					if !cancelled {
						processError(writer, err)
					}
				}))
			return b
		}))
}

func (h *httpServer) generateLinks(prefix string) maps.String {
	var depth int
	if prefix != "/" {
		tempEndpointSplitter := new(SplitEndpoint)
		tempEndpointSplitter.SetEndpoint(prefix, '/')
		depth = tempEndpointSplitter.NumParts()
	}
	prefix = strings.TrimLeft(prefix, "/")
	var links maps.String
	if depth >= len(h.endpoints.groupEndpoints) {
		return nil
	}
	for key := range h.endpoints.groupEndpoints[depth] {
		if strings.HasPrefix(key, prefix) && key != prefix {
			key = strings.Replace(key, ".", "/", -1)
			//parts := strings.Split(urlPathName, ".")
			//links.Add("/"+key, fmt.Sprintf("http://%v/%v/%v/%v", h.server.Addr, services.DefaultHttpServerName, key, parts[len(parts) - 1]))
			links.Add("/"+key, fmt.Sprintf("http://%v/doc/%v/", h.server.Addr, key))
		}
	}
	return links
}

func (h *httpServer) generateMethods(url string) (map[string]map[string]endpointInfo, error) {
	data := make(map[string]map[string]endpointInfo)
	path, method, methods := h.endpoints.FindEndpointMethods(url, '/')
	if len(path) > 0 {
		return nil, basicerrors.BadParameter
	}
	data[url] = make(map[string]endpointInfo)
	if len(method) > 0 {
		return nil, basicerrors.NotFound
	}
	for methodName, endpoints := range methods {
		for httpMethodName, endpoint := range endpoints {
			if endpoint.Dest != nil {
				data[url][fmt.Sprintf("[%v] %v", httpMethodName, methodName)] = *endpoint
			}
		}
	}
	return data, nil
}

func (h *httpServer) generateDoc(cmd *urlPath) (replies.StringReply, error) {
	realUrl := strings.Replace(cmd.id, ".", "/", -1)
	realUrl = "/" + realUrl
	methods, err := h.generateMethods(realUrl)
	if err != nil {
		return "", err
	}

	links := h.generateLinks(realUrl)

	result := new(strings.Builder)
	err = executeMethodsTemplate(result, methods)
	if err != nil {
		return "", err
	}
	if links != nil {
		err = executeLinksTemplate(result, links)
		if err != nil {
			return "", err
		}
	}
	return replies.StringReply(result.String()), nil
}

func init() {
	var sessionDuration string
	defaultHttpServerParams := flags.HostPort{Port: 8882}
	defaultDir, _ := os.Getwd()
	starter.SetCreator(DefaultHttpServerName, func(s *actors.Actor, name string) (actors.ActorService, error) {
		if len(sessionDuration) > 0 {
			duration, err := time.ParseDuration(sessionDuration)
			if err != nil {
				return nil, err
			}
			defaultConfig.SessionDuration = duration
		}
		server := NewHttpServer(defaultHttpServerParams, DefaultHttpServerName, defaultDir, &defaultConfig)
		return s.System().Spawn(server), nil
	})

	starter.SetFlagInitializer(DefaultHttpServerName, func() {
		defaultHttpServerParams.RegisterFlagsWithDescriptions(
			"http",
			"listen to http requests on this hostname/ip address",
			"listen to http requests on this port")
		flags.StringFlag(&defaultDir, "endpoints", "saved endpoints directory")
		flags.StringFlag(&defaultConfig.JwtKey, "jwt", "jwt secret key")
		flags.StringFlag(&sessionDuration, "session-duration", "duration of session (format: 5d4h3m2s")
		flags.Float64Flag(&defaultConfig.SessionResetPercent, "session-reset", "remaining session percentage when httpServer automatically refresh session cookie")
	})
}
