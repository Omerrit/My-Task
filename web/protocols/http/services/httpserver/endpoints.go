package httpserver

import (
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/basicerrors"
	"gerrit-share.lan/go/debug"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"gerrit-share.lan/go/interfaces"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/metadata"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/sets"
	"strings"
)

const defaultHttpMethod = "get"
const (
	endpointInfoName = packageName + ".endpointinfo"
	endpointsName    = packageName + ".endpoints"
)

type endpointInfo struct {
	Dest             actors.ActorService
	CommandGenerator inspectables.Creator
	CommandMetaData  *metadata.Metadata
	ResultInfo       *metadata.Metadata
	OriginalName     string
	Resource         string
	Method           string
	HttpMethod       string
	changed          bool
}

func (e *endpointInfo) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(endpointInfoName, "endpoint information")
	{
		objectInspector.String(&e.HttpMethod, "http", false, "http method")
		objectInspector.String(&e.Method, "method", false, "method")
		objectInspector.String(&e.Resource, "resource", false, "resource")
	}
	objectInspector.End()
}

func (e *endpointInfo) update(resource, method, httpMethod string) {
	e.Method = method
	e.HttpMethod = httpMethod
	e.Resource = resource
	e.changed = true
}

func init() {
	inspectables.Register(endpointInfoName, func() inspect.Inspectable { return new(endpointInfo) })
}

type endpointInfoByName map[string]*endpointInfo

func (e *endpointInfoByName) init(name, resource, method, httpMethod string) *endpointInfo {
	if *e == nil {
		*e = make(endpointInfoByName, 1)
	}
	info := &endpointInfo{
		OriginalName: name,
		Resource:     resource,
		Method:       method,
		HttpMethod:   httpMethod,
		changed:      true,
	}
	(*e)[name] = info
	return info
}

func (e *endpointInfoByName) add(name string, dest actors.ActorService, creator inspectables.Creator,
	command *metadata.Metadata, result *metadata.Metadata) *endpointInfo {
	if *e == nil {
		*e = make(map[string]*endpointInfo, 1)
	}
	info := &endpointInfo{
		Dest:             dest,
		CommandGenerator: creator,
		CommandMetaData:  command,
		ResultInfo:       result,
		OriginalName:     name,
	}
	(*e)[name] = info
	return info
}

func (e *endpointInfoByName) removeByDestination(service interfaces.ClosableService) {
	actor, ok := service.(actors.ActorService)
	if !ok {
		return
	}
	for name, info := range *e {
		if info.Dest == actor {
			delete(*e, name)
		}
	}
}

func (e *endpointInfoByName) remove(name string) {
	delete(*e, name)
}

func (e *endpointInfoByName) Visit(v actors.ResponseVisitor) {
	v.Reply(e)
}

func (e *endpointInfoByName) Inspect(i *inspect.GenericInspector) {
	mapInspector := i.Map(endpointsName, endpointInfoName, "list of available endpoints")
	if mapInspector.IsReading() {
		return
	}
	var lengthToWrite int
	for _, value := range *e {
		if value.changed {
			lengthToWrite++
		}
	}
	mapInspector.SetLength(lengthToWrite)
	for key, value := range *e {
		if value.changed {
			valI := mapInspector.WriteValue(key)
			value.Inspect(valI)
		}
	}
	mapInspector.End()
}

func init() {
	inspectables.Register(endpointsName, func() inspect.Inspectable { return new(endpointInfoByName) })
}

type endpointByHttpMethod map[string]*endpointInfo

func (e *endpointByHttpMethod) add(httpMethod string, info *endpointInfo) {
	if *e == nil {
		*e = make(endpointByHttpMethod, 1)
	}
	(*e)[httpMethod] = info
}

func (e *endpointByHttpMethod) remove(httpMethod string) {
	delete(*e, httpMethod)
}

func (e *endpointByHttpMethod) isEmpty() bool {
	return len(*e) == 0
}

type endpointInfoByMethod map[string]endpointByHttpMethod

func (e *endpointInfoByMethod) add(method string, httpMethod string, info *endpointInfo) {
	if *e == nil {
		*e = make(endpointInfoByMethod, 1)
	}
	infos := (*e)[method]
	infos.add(httpMethod, info)
	(*e)[method] = infos
}

func (e *endpointInfoByMethod) remove(method string, httpMethod string) {
	v, ok := (*e)[method]
	if ok {
		v.remove(httpMethod)
		if v.isEmpty() {
			delete(*e, method)
		}
	}
}

func (e *endpointInfoByMethod) isEmpty() bool {
	return len(*e) == 0
}

//endpoint name->endpoints
type simpleEndpoints map[string]endpointInfoByMethod

func (s *simpleEndpoints) add(name string, method string, httpMethod string, info *endpointInfo) {

	if *s == nil {
		*s = make(simpleEndpoints, 1)
	}
	infos := (*s)[name]
	infos.add(method, httpMethod, info)
	(*s)[name] = infos
}

func (s *simpleEndpoints) remove(name string, method string, httpMethod string) {
	infos, ok := (*s)[name]
	if ok {
		infos.remove(method, httpMethod)
		if infos.isEmpty() {
			delete(*s, name)
		}
	}
}

//endpoint depth->endpoints
type groupEndpoints []simpleEndpoints

func (e *groupEndpoints) add(name string, method string, depth int, httpMethod string, info *endpointInfo) {
	if len(*e) <= depth {
		*e = append(*e, make(groupEndpoints, depth+1-len(*e))...)
	}
	(*e)[depth].add(name, method, httpMethod, info)
}

func (e *groupEndpoints) remove(name string, method string, httpMethod string) {
	depth := getResourceDepth(name)
	if depth >= len(*e) {
		return
	}
	(*e)[depth].remove(name, method, httpMethod)
}

func (e *groupEndpoints) getInfo(name string, method string, depth int, httpMethod string) *endpointInfo {
	if depth >= len(*e) {
		return nil
	}
	return (*e)[depth][name][method][httpMethod]
}

func (e *groupEndpoints) isEmpty() bool {
	return len(*e) == 0
}

func (e *groupEndpoints) size() int {
	return len(*e)
}

type HttpRestEndpoints struct {
	groupEndpoints       groupEndpoints
	destinationEndpoints sets.StringsForServices
	info                 endpointInfoByName
}

func (h *HttpRestEndpoints) Add(name string, service actors.ActorService, generator inspectables.Creator,
	resultSample *metadata.Metadata, commandMetadata *metadata.Metadata) {
	h.destinationEndpoints.Add(service, name)
	if info, ok := h.info[name]; ok {
		info.CommandMetaData = commandMetadata
		info.Dest = service
		info.ResultInfo = resultSample
		info.CommandGenerator = generator
		return
	}

	parsed := ParseEndpoint(name)
	debug.Printf("%#v\n", parsed)
	info := h.info.add(name, service, generator, commandMetadata, resultSample)
	h.groupEndpoints.add(parsed.Name, parsed.Method, parsed.Depth, defaultHttpMethod, info)
}

func (h *HttpRestEndpoints) Load(resource string, method string, httpMethod string, originalName string) {
	debug.Printf("%v loaded. resource: %v, method: %v, http method: %v\n", originalName, resource, method, httpMethod)
	info := h.info.init(originalName, resource, method, httpMethod)
	h.groupEndpoints.add(resource, method, getResourceDepth(resource), httpMethod, info)
}

func (h *HttpRestEndpoints) Edit(resource, method, httpMethod string, newResource, newMethod, newHttpMethod string) error {
	depth := getResourceDepth(resource)
	if depth > len(h.groupEndpoints)-1 {
		return basicerrors.NotFound
	}
	info, ok := h.groupEndpoints[depth][resource][method][httpMethod]
	if !ok {
		return basicerrors.NotFound
	}
	info.update(newResource, newMethod, newHttpMethod)
	h.groupEndpoints.remove(resource, method, httpMethod)
	h.groupEndpoints.add(newResource, newMethod, getResourceDepth(newResource), newHttpMethod, info)
	return nil
}

func (h *HttpRestEndpoints) EditByName(name string, newResource, newMethod, newHttpMethod string) error {
	info, ok := h.info[name]
	if !ok {
		return basicerrors.NotFound
	}
	h.groupEndpoints.remove(info.Resource, info.Method, info.HttpMethod)
	info.update(newResource, newMethod, newHttpMethod)
	h.groupEndpoints.add(newResource, newMethod, getResourceDepth(newResource), newHttpMethod, info)
	return nil
}

func (h *HttpRestEndpoints) Remove(name string, httpMethod string) {
	parsed := ParseEndpoint(name)
	service := h.groupEndpoints.getInfo(parsed.Name, parsed.Method, parsed.Depth, httpMethod).Dest
	h.groupEndpoints.remove(parsed.Name, parsed.Method, httpMethod)
	h.destinationEndpoints.Remove(service, name)
}

func (h *HttpRestEndpoints) RemoveByDestination(destination actors.ActorService) {
	for name := range h.destinationEndpoints[destination] {
		info := h.info[name]
		if info.changed {
			h.groupEndpoints.remove(info.Resource, info.Method, info.HttpMethod)
		} else {
			parsed := ParseEndpoint(name)
			h.groupEndpoints.remove(parsed.Name, parsed.Method, defaultHttpMethod)
		}
		h.info.remove(name)
	}
	h.destinationEndpoints.RemoveByDestination(destination)
}

func (h *HttpRestEndpoints) findEndpointMethods1Part(fullPath []byte) (path []byte, method []byte, endpoints endpointInfoByMethod) {
	debug.Println("1part:", string(fullPath))
	result, ok := h.groupEndpoints[0][string(fullPath)]
	if !ok {
		result, ok := h.groupEndpoints[0][""]
		if !ok {
			return nil, nil, nil
		}
		return nil, fullPath, result
	}
	return nil, nil, result
}

func (h *HttpRestEndpoints) findEndpointMethods2Part(fullPath []byte, possibleEndpoint []byte,
	possibleMethod []byte) (path []byte, method []byte, endpoints endpointInfoByMethod) {

	debug.Println("2part:", string(possibleEndpoint), string(possibleMethod))
	result, ok := h.groupEndpoints[0][string(possibleEndpoint)]
	if ok {
		return nil, possibleMethod, result
	}
	return possibleEndpoint, possibleMethod, h.groupEndpoints[0][""]
}

func (h *HttpRestEndpoints) FindEndpointMethods(endpoint string, separator byte) (path []byte, method []byte, methods endpointInfoByMethod) {
	debug.Println("findMethods:", endpoint)
	split := new(SplitEndpoint)
	split.SetEndpoint(endpoint, separator)
	if h.groupEndpoints.isEmpty() {
		return nil, nil, nil
	}
	if split.NumParts() <= h.groupEndpoints.size() {
		debug.Println("full path:", string(split.FullPath()))
		result, ok := h.groupEndpoints[split.NumParts()-1][string(split.FullPath())]
		if ok {
			return nil, nil, result
		}
	}

	switch split.NumParts() {
	case 1:
		return h.findEndpointMethods1Part(split.FullPath())
	case 2:
		return h.findEndpointMethods2Part(split.FullPath(), split.CurrentPath(), split.Tail())
	}
	candidateMethod := split.Tail()
	split.Cut(h.groupEndpoints.size() + 1)
	for ; split.NumParts() > 1; split.Advance() {
		endpoint := split.CurrentPath()
		methods = h.groupEndpoints[split.NumParts()-2][string(endpoint)]
		if methods != nil {
			return split.CuttedFullPath(len(endpoint)+1, len(candidateMethod)+1), candidateMethod, methods
		}
	}
	return split.CuttedFullPath(0, len(candidateMethod)+1), candidateMethod, h.groupEndpoints[0][""]
}

func (h *HttpRestEndpoints) getEndpointByOriginalName(name string) (*endpointInfo, error) {
	endpoint := &endpointInfo{}
	index := strings.Index(name, string(hierarchySeparator))
	if index == -1 {
		return endpoint, ErrHttpMethodNotAllowed
	}
	httpMethod := strings.ToLower(name[:index])
	info := ParseEndpoint(name[index+1:])
	if info.Depth+1 > len(h.groupEndpoints) {
		return endpoint, ErrResourceNotFound
	}
	methods, ok := h.groupEndpoints[info.Depth][info.Name]
	if !ok {
		return endpoint, ErrResourceNotFound
	}
	httpMethods, ok := methods[info.Method]
	if !ok {
		return endpoint, ErrMethodNotAllowed
	}
	endpoint, ok = httpMethods[httpMethod]
	if !ok {
		return endpoint, ErrHttpMethodNotAllowed
	}
	return endpoint, nil
}

func (h *HttpRestEndpoints) Inspect(i *inspect.GenericInspector) {
	mapInspector := i.Map("", "", "endpoint list")
	if mapInspector.IsReading() {
		length := mapInspector.GetLength()
		for i := 0; i < length; i++ {
			key := mapInspector.NextKey()
			genericInspector := mapInspector.ReadValue()
			objectInspector := genericInspector.Object("", "")
			var resource, method, httpMethod string
			objectInspector.String(&resource, "resource", false, "")
			objectInspector.String(&method, "method", false, "")
			objectInspector.String(&httpMethod, "http", false, "")
			objectInspector.End()
			h.Load(resource, method, httpMethod, key)
		}
	} else {
		mapInspector.SetLength(len((*h).groupEndpoints))
		for _, endpoints := range (*h).groupEndpoints {
			for resource, endpoints := range endpoints {
				for method, endpoints := range endpoints {
					for httpMethod, info := range endpoints {
						objectInspector := mapInspector.WriteValue(info.OriginalName).Object("", "position object")
						objectInspector.String(&resource, "resource", false, "")
						objectInspector.String(&method, "method", false, "")
						objectInspector.String(&httpMethod, "http", false, "")
						objectInspector.End()
					}
				}
			}
		}
	}
	mapInspector.End()
}

func (h *HttpRestEndpoints) print() {
	for k, v := range h.groupEndpoints {
		fmt.Println(k)
		for k, v := range v {
			fmt.Println("\t", k)
			for k, v := range v {
				fmt.Println("\t\t", k)
				for k, v := range v {
					fmt.Println("\t\t\t", k, ":", v.OriginalName)
				}
			}
		}
	}
}
