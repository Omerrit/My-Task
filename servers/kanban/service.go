package kanban

import (
	"bytes"
	"context"
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/replies"
	"gerrit-share.lan/go/actors/services/shutdownactor"
	"gerrit-share.lan/go/actors/starter"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/json/fromjson"
	"gerrit-share.lan/go/servers/kanban/internal/endpoints"
	"gerrit-share.lan/go/servers/kanban/internal/ids"
	"gerrit-share.lan/go/servers/kanban/internal/kafka"
	"gerrit-share.lan/go/servers/kanban/internal/utils"
	"gerrit-share.lan/go/utils/flags"
	"gerrit-share.lan/go/utils/sets"
	"gerrit-share.lan/go/web/protocols/http/serializers"
	"github.com/Shopify/sarama"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const KanbanName = "kanban"

const (
	streamPath    = "/stream"
	loginPath     = "/login"
	intPath       = "/int"
	floatPath     = "/float"
	boolPath      = "/bool"
	stringPath    = "/string"
	dateTimePath  = "/datetime"
	datePath      = "/date"
	timePath      = "/time"
	userIdPath    = "/userid"
	newIdPath     = "/newid"
	deleteIdPath  = "/deleteid"
	reserveIdPath = "/reserveid"
	saveFilePath  = "/savefile"
)

type kanban struct {
	actors.Actor
	server        http.Server
	name          string
	httpActor     actors.ActorService
	broadcaster   actors.StateBroadcaster
	messages      kafka.MessagesStream
	config        *sarama.Config
	brokers       []string
	keyValueTopic string
	historyTopic  string
	producer      sarama.SyncProducer
	usersStorage  utils.UsersStorage
	endpoints     endpoints.Endpoints
	ids           ids.TypedIds
	topicKeys     sets.String
	configLength  int
	eofReached    bool
	copyToHistory bool
	workDir       string
}

func newKanban(name string, system *actors.System, dir string, httpHostPort flags.HostPort, kafkaBrokers []string, keyValueTopic string, historyTopic string) (actors.ActorService, error) {
	if len(keyValueTopic) == 0 || len(historyTopic) == 0 {
		return nil, fmt.Errorf("empty topic name detected")
	}
	if keyValueTopic == historyTopic {
		return nil, fmt.Errorf("same topic provided for history and key value")
	}
	server := new(kanban)
	server.name = name
	server.keyValueTopic = keyValueTopic
	server.historyTopic = historyTopic
	server.configLength = len(typeMessages) + len(propMessages)
	server.brokers = kafkaBrokers
	server.workDir = dir
	server.config = kafka.NewConfig()
	server.generateEndpoints()
	server.server = http.Server{
		Addr:    httpHostPort.String(),
		Handler: server}
	err := kafka.CheckTopic(server.keyValueTopic, server.brokers, server.config, kafka.CompactTopicEntries())
	if err != nil {
		return nil, err
	}
	err = kafka.CheckTopic(server.historyTopic, server.brokers, server.config, kafka.HistoryTopicEntries())
	if err != nil {
		return nil, err
	}
	server.copyToHistory, err = server.hasToCopyToHistory()
	if err != nil {
		return nil, err
	}
	server.producer, err = sarama.NewSyncProducer(server.brokers, server.config)
	if err != nil {
		return nil, err
	}
	return system.Spawn(server), nil
}

func (k *kanban) MakeBehaviour() actors.Behaviour {
	log.Println(k.name, "started")
	var starterHandle starter.Handle
	starterHandle.Acquire(k, starterHandle.DependOn, k.Quit)

	k.broadcaster = actors.NewBroadcaster(&k.messages)
	k.broadcaster.CloseWhenActorCloses()

	var behaviour actors.Behaviour
	behaviour.Name = k.name
	behaviour.AddCommand(new(subscribe), func(cmd interface{}) (actors.Response, error) {
		output := k.broadcaster.AddOutput()
		source := output.DataSource().(*kafka.MessageDataSource)
		startId := cmd.(*subscribe).startId
		lastId, _ := k.messages.GetLatestState()
		if startId > lastId {
			source.Init(startId - 1)
			log.Println("consumer not initialized")
			k.InitStreamOutput(output, cmd.(*subscribe))
			return nil, nil
		}
		if startId == 0 {
			k.consumeFromKafka(sarama.OffsetOldest, source.Init(-1))
			log.Println("consumer initialized")
		} else {
			k.consumeFromKafka(startId, source.Init(startId-1))
			log.Println("consumer initialized")
		}
		k.InitStreamOutput(output, cmd.(*subscribe))
		return nil, nil
	})
	behaviour.AddCommand(new(login), func(cmd interface{}) (actors.Response, error) {
		loginCmd := cmd.(*login)
		return replies.Bool(k.usersStorage.AreCredentialsValid(loginCmd.userName, loginCmd.password)), nil
	}).ResultBool()
	behaviour.AddCommand(new(newId), func(cmd interface{}) (actors.Response, error) {
		newIdCmd := cmd.(*newId)
		return replies.String(k.ids.AcquireNewId(newIdCmd.objectType, newIdCmd.id)), nil
	}).ResultString()
	behaviour.AddCommand(new(deleteId), func(cmd interface{}) (actors.Response, error) {
		deleteCmd := cmd.(*deleteId)
		return nil, k.ids.DeleteId(deleteCmd.objectType, deleteCmd.id)
	})
	behaviour.AddCommand(new(isIdRegistered), func(cmd interface{}) (actors.Response, error) {
		isRegisteredCmd := cmd.(*isIdRegistered)
		return replies.Bool(k.ids.IsRegistered(isRegisteredCmd.objectType, isRegisteredCmd.id)), nil
	}).ResultBool()
	behaviour.AddCommand(new(reserveId), func(cmd interface{}) (actors.Response, error) {
		reserveIdCmd := cmd.(*reserveId)
		return nil, k.ids.RestoreId(reserveIdCmd.objectType, reserveIdCmd.id)
	})
	behaviour.AddCommand(new(saveMsgsToKafka), func(cmd interface{}) (actors.Response, error) {
		saveCmd := cmd.(*saveMsgsToKafka)
		return nil, k.saveMessagesToKafka(saveCmd)
	})
	k.SetPanicProcessor(k.onPanic)

	err := k.consumeFromKafka(0, k.writeMessage)
	if err != nil {
		k.Quit(err)
		return behaviour
	}
	return behaviour
}

func (k *kanban) Shutdown() error {
	err := k.server.Shutdown(context.Background())
	if err != nil {
		log.Println("error while shutting of http server down:", err)
	}
	err = k.producer.Close()
	if err != nil {
		log.Println("error while shutting of kafka producer:", err)
	}
	log.Println(k.name, "shut down")
	return nil
}

func (k *kanban) onPanic(err errors.StackTraceError) {
	log.Println("panic:", err, err.StackTrace())
	k.Quit(err)
}

func (k *kanban) generateEndpoints() {
	editCreator := func() inspect.Inspectable {
		return &saveMsgToKafka{}
	}
	k.endpoints.Add(intPath, k.edit, editCreator)
	k.endpoints.Add(floatPath, k.edit, editCreator)
	k.endpoints.Add(boolPath, k.edit, editCreator)
	k.endpoints.Add(stringPath, k.edit, editCreator)
	k.endpoints.Add(dateTimePath, k.edit, editCreator)
	k.endpoints.Add(datePath, k.edit, editCreator)
	k.endpoints.Add(timePath, k.edit, editCreator)
	k.endpoints.Add(userIdPath, k.edit, editCreator)
	k.endpoints.Add(streamPath, k.openStream, nil)
	k.endpoints.Add(loginPath, k.login, func() inspect.Inspectable {
		return &login{}
	})
	k.endpoints.Add(newIdPath, k.newId, func() inspect.Inspectable {
		return &newId{}
	})
	k.endpoints.Add(deleteIdPath, k.deleteId, func() inspect.Inspectable {
		return &deleteId{}
	})
	k.endpoints.Add(reserveIdPath, k.reserveId, func() inspect.Inspectable {
		return &reserveId{}
	})
	k.endpoints.Add(saveFilePath, k.saveFile, func() inspect.Inspectable {
		return &saveFile{}
	})
}

func (k *kanban) consumeFromKafka(offset int64, processor func(*kafka.Message) (error, bool)) error {
	consumer, err := kafka.NewConsumer("kafka_consumer", k.System(), k.config, k.brokers, k.keyValueTopic, 0, offset)
	if err != nil {
		return err
	}

	stopConsuming := false

	kafkaInput := actors.NewSimpleCallbackStreamInput(func(data inspect.Inspectable) error {
		msgs := data.(*kafka.Messages)
		for _, msg := range *msgs {
			err, stopConsuming = processor(msg)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(base *actors.StreamInputBase) {
		if stopConsuming {
			log.Println("closing input")
			base.Acknowledge()
			base.Close()
		} else {
			base.RequestData(new(kafka.Messages), 10)
		}
	})

	kafkaInput.CloseWhenActorCloses()
	//kafkaInput.OnClose(k.Quit)
	k.RequestStream(kafkaInput, consumer, &kafka.Subscribe{}, k.Quit)
	return nil
}

func (k *kanban) runHttp() {
	k.httpActor = k.System().RunAsyncSimpleNamed("http server", func() error {
		log.Println("listen and serve started")
		fmt.Println(k.server.ListenAndServe())
		log.Println("listen and serve shutdown")
		return nil
	})
	k.DependOn(k.httpActor)
}

func (k *kanban) serveStatic(urlPath string, writer http.ResponseWriter) {
	filePath := strings.TrimPrefix(urlPath, "/static/")
	file, err := os.Open(path.Join(k.workDir, filePath))
	defer file.Close()
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("file not found"))
		return
	}
	buff := bytes.Buffer{}
	_, err = buff.ReadFrom(file)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("failed to read from file"))
		return
	}
	writer.Write(buff.Bytes())
}

func (k *kanban) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	if request.Method == "OPTIONS" {
		writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		writer.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimSuffix(request.URL.Path, "/")
	if strings.HasPrefix(path, "/static/") {
		k.serveStatic(path, writer)
		return
	}

	endpoint, ok := k.endpoints[path]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	switch contentType := request.Header.Get("Content-Type"); contentType {
	case "application/json":
		k.processJsonRequest(request, endpoint, writer)
	default:
		if strings.HasPrefix(contentType, "multipart/form-data") {
			k.processFormData(request, endpoint, writer)
			return
		}

		if endpoint.Creator() != nil {
			writer.WriteHeader(http.StatusUnsupportedMediaType)
			writer.Write([]byte("unsupported content type"))
			return
		}
		endpoint.Handler()(request.Context(), nil, request.Header, nil, writer)
	}
}

func (k *kanban) processFormData(request *http.Request, endpoint endpoints.Endpoint, writer http.ResponseWriter) {
	values, file, err := utils.ParseMultiForm(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	command := endpoint.Creator()()
	deserializer := &serializers.FromUrl{Values: values}
	i := inspect.NewGenericInspector(deserializer)
	command.Inspect(i)
	if i.GetError() != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(i.GetError().Error()))
		return
	}

	endpoint.Handler()(request.Context(), command, request.Header, file, writer)
}

func (k *kanban) processJsonRequest(request *http.Request, endpoint endpoints.Endpoint, writer http.ResponseWriter) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	parser := fromjson.NewInspector(body, 0)
	inspector := inspect.NewGenericInspector(parser)
	var command inspect.Inspectable
	if endpoint.Creator() != nil {
		command = endpoint.Creator()()
		command.Inspect(inspector)
		if inspector.GetError() != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(inspector.GetError().Error()))
			return
		}
	}
	endpoint.Handler()(request.Context(), command, request.Header, nil, writer)
}

func (k *kanban) saveFile(_ context.Context, command inspect.Inspectable, _ http.Header, file *multipart.Part, writer http.ResponseWriter) {
	onRequestError := func(err error) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("failed to save file: " + err.Error()))
	}

	var fileType string
	if val, ok := file.Header["Content-Type"]; !ok && len(val) > 0 {
		onRequestError(fmt.Errorf("file type not provided"))
		return
	}
	fileType = file.Header["Content-Type"][0]

	saveFileCommand := command.(*saveFile)
	directory := path.Join(k.workDir, path.Join(strings.Split(saveFileCommand.id, ids.IdSeparator)...))
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		onRequestError(err)
		return
	}

	dirFiles, _ := ioutil.ReadDir(directory)
	intFileSrcName := 1
	for _, val := range dirFiles {
		intFileName, err := strconv.Atoi(val.Name())
		if err != nil {
			continue
		}
		intFileName++
		if intFileName > intFileSrcName {
			intFileSrcName = intFileName
		}
	}

	filePath := path.Join(path.Join(strings.Split(saveFileCommand.id, ids.IdSeparator)...), strconv.Itoa(intFileSrcName))
	fileToSave, err := os.OpenFile(path.Join(k.workDir, filePath), os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		onRequestError(err)
		return
	}
	defer file.Close()

	_, err = io.Copy(fileToSave, file)
	if err != nil {
		onRequestError(err)
		return
	}

	var messages saveMsgsToKafka
	const fileTypeName = "file"

	messages = append(messages, saveMsgToKafka{
		utils.NewValue(saveFileCommand.user, time.Now().Unix(), file.FileName(), "").ToJson(),
		utils.NewParsedKey(fileTypeName, saveFileCommand.id, "name").ToString(),
	})
	messages = append(messages, saveMsgToKafka{
		utils.NewValue(saveFileCommand.user, time.Now().Unix(), fileType, "").ToJson(),
		utils.NewParsedKey(fileTypeName, saveFileCommand.id, "type").ToString(),
	})
	messages = append(messages, saveMsgToKafka{
		utils.NewValue(saveFileCommand.user, time.Now().Unix(), filePath, "").ToJson(),
		utils.NewParsedKey(fileTypeName, saveFileCommand.id, "src").ToString(),
	})
	messages = append(messages, saveMsgToKafka{
		utils.NewValue(saveFileCommand.user, time.Now().Unix(), strconv.Itoa(saveFileCommand.lastModified), "").ToJson(),
		utils.NewParsedKey(fileTypeName, saveFileCommand.id, "last_modified").ToString(),
	})

	k.System().Become(actors.NewSimpleActor(func(actor *actors.Actor) actors.Behaviour {
		behaviour := actors.Behaviour{Name: "client saveFile handler"}
		actor.SendRequest(k.Service(), &messages,
			actors.OnReply(func(reply interface{}) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte("ok"))
			}).OnError(func(err error) {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte("failed to save file data to kafka: " + err.Error()))
			}))
		return behaviour
	}))
}

func (k *kanban) saveMessagesToKafka(messages *saveMsgsToKafka) error {
	for _, msg := range *messages {
		err := k.sendMessageToKafka(msg.key, msg.value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *kanban) reserveId(context context.Context, command inspect.Inspectable, _ http.Header, _ *multipart.Part, writer http.ResponseWriter) {
	reserveIdCmd := command.(*reserveId)
	k.System().Become(actors.NewSimpleActor(func(actor *actors.Actor) actors.Behaviour {
		behaviour := actors.Behaviour{Name: "client reserveId handler"}
		actor.SendRequest(k.Service(), reserveIdCmd,
			actors.OnReply(func(reply interface{}) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte("ok"))
			}).OnError(func(err error) {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte("failed to reserve id: " + err.Error()))
			}))
		return behaviour
	}))
}

func (k *kanban) deleteId(context context.Context, command inspect.Inspectable, _ http.Header, _ *multipart.Part, writer http.ResponseWriter) {
	deleteIdCmd := command.(*deleteId)
	k.System().Become(actors.NewSimpleActor(func(actor *actors.Actor) actors.Behaviour {
		behaviour := actors.Behaviour{Name: "client deleteId handler"}
		actor.SendRequest(k.Service(), deleteIdCmd,
			actors.OnReply(func(reply interface{}) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte("ok"))
			}).OnError(func(err error) {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte("failed to delete id: " + err.Error()))
			}))
		return behaviour
	}))
}

func (k *kanban) newId(context context.Context, command inspect.Inspectable, _ http.Header, _ *multipart.Part, writer http.ResponseWriter) {
	newIdCmd := command.(*newId)
	k.System().Become(actors.NewSimpleActor(func(actor *actors.Actor) actors.Behaviour {
		behaviour := actors.Behaviour{Name: "client newId handler"}
		actor.SendRequest(k.Service(), newIdCmd,
			actors.OnReply(func(reply interface{}) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(reply.(string)))
			}).OnError(func(err error) {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("failed to acquire id: " + err.Error()))
			}))
		return behaviour
	}))
}

func (k *kanban) login(context context.Context, command inspect.Inspectable, _ http.Header, _ *multipart.Part, writer http.ResponseWriter) {
	loginCmd := command.(*login)
	k.System().Become(actors.NewSimpleActor(func(actor *actors.Actor) actors.Behaviour {
		behaviour := actors.Behaviour{Name: "client login handler"}
		actor.SendRequest(k.Service(), loginCmd,
			actors.OnReply(func(reply interface{}) {
				ok := reply.(bool)
				if ok {
					writer.WriteHeader(http.StatusOK)
					writer.Write([]byte("ok"))
				} else {
					writer.WriteHeader(http.StatusBadRequest)
					writer.Write([]byte("incorrect password or user name"))
				}
			}).OnError(func(err error) {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("something wrong with service:" + err.Error()))
			}))
		return behaviour
	}))
}

func (k *kanban) isTopicValid() bool {
	if len(k.topicKeys) > k.configLength {
		return len(typeMessages) == 0 && len(propMessages) == 0
	}
	return true
}

func (k *kanban) isConfigMessage(message *kafka.Message) bool {
	_, isType := typeMessages[string(message.Key)]
	_, isProp := propMessages[string(message.Key)]

	if !isType && !isProp {
		return false
	}
	if isProp {
		propMessages.Delete(string(message.Key))
	}
	if isType {
		typeMessages.Delete(string(message.Key))
	}
	return true
}

func (k *kanban) sendMessageToKafka(key string, valueAsString string) error {
	var value sarama.Encoder
	if len(valueAsString) > 0 {
		value = sarama.StringEncoder(valueAsString)
	}
	message := &sarama.ProducerMessage{
		Topic:     k.keyValueTopic,
		Key:       sarama.StringEncoder(key),
		Value:     value,
		Headers:   nil,
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
		Timestamp: time.Now(),
	}
	_, _, err := k.producer.SendMessage(message)
	if err != nil {
		return err
	}
	message.Topic = k.historyTopic
	_, _, err = k.producer.SendMessage(message)
	return err
}

func (k *kanban) sendRestOfConfig() error {
	failed := func(err error) error {
		return fmt.Errorf("failed to write message to key value topic: %s", err.Error())
	}
	for key, value := range typeMessages {
		err := k.sendMessageToKafka(key, value)
		if err != nil {
			return failed(err)
		}
	}
	for key, value := range propMessages {
		err := k.sendMessageToKafka(key, value)
		if err != nil {
			return failed(err)
		}
	}
	return nil
}

func (k *kanban) hasToCopyToHistory() (bool, error) {
	failed := func(err error) (bool, error) {
		return false, fmt.Errorf("failed to write history: %s", err.Error())
	}
	client, err := sarama.NewClient(k.brokers, k.config)
	if err != nil {
		return failed(err)
	}
	defer client.Close()
	offsetNewest, err := client.GetOffset(k.historyTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return failed(err)
	}
	offsetOldest, err := client.GetOffset(k.historyTopic, 0, sarama.OffsetOldest)
	if err != nil {
		return failed(err)
	}
	return offsetOldest == offsetNewest, nil
}

func (k *kanban) writeMessage(command *kafka.Message) (error, bool) {
	handlerError := func(err error) (error, bool) {
		log.Println(err)
		k.Quit(err)
		return err, false
	}
	if command.Offset == kafka.OffsetUninitialized {
		k.eofReached = true
		if !k.isTopicValid() {
			return handlerError(fmt.Errorf("shitty topic detected! use empty or properly configured topic next time"))
		}
		k.topicKeys.Clear()
		err := k.sendRestOfConfig()
		if err != nil {
			return handlerError(err)
		}
		k.runHttp()
		return nil, false
	}

	if !k.eofReached && len(k.topicKeys) < k.configLength {
		if !k.isConfigMessage(command) {
			err := fmt.Errorf("shitty topic detected! use empty or properly configured topic next time")
			log.Println(err)
			k.Quit(err)
			return err, false
		}
		k.topicKeys.Add(string(command.Key))
	}

	k.messages.Add(command)
	k.broadcaster.NewDataAvailable()

	if !k.eofReached {
		parsedKey, err := utils.ParseKey(string(command.Key))
		if err == nil {
			k.ids.RestoreId(parsedKey.Type, parsedKey.Id)
		}
		if k.copyToHistory {
			_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
				Topic:     k.historyTopic,
				Key:       sarama.ByteEncoder(command.Key),
				Value:     sarama.ByteEncoder(command.Value),
				Headers:   nil,
				Metadata:  nil,
				Offset:    0,
				Partition: 0,
				Timestamp: time.Now(),
			})
			if err != nil {
				return handlerError(err)
			}
		}
	}

	key, err := utils.ParseKey(string(command.Key))
	if err != nil {
		return nil, false
	}
	switch key.Type {
	case utils.UserTypeName:
		k.usersStorage.ProcessUserProp(key, command.Value)
	}
	return nil, false
}

func (k *kanban) edit(context context.Context, command inspect.Inspectable, _ http.Header, _ *multipart.Part, writer http.ResponseWriter) {
	msgCommand := command.(*saveMsgToKafka)
	parsedKey, err := utils.ParseKey(msgCommand.key)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("incorrect key format"))
		return
	}
	k.System().Become(actors.NewSimpleActor(func(actor *actors.Actor) actors.Behaviour {
		behaviour := actors.Behaviour{Name: "client edit handler"}
		actor.SendRequest(k.Service(), &isIdRegistered{parsedKey.Type, parsedKey.Id},
			actors.OnReply(func(reply interface{}) {
				ok := reply.(bool)
				if ok {
					var msgs saveMsgsToKafka
					msgs = append(msgs, *msgCommand)
					actor.SendRequest(k.Service(), &msgs,
						actors.OnReply(func(reply interface{}) {
							writer.WriteHeader(http.StatusOK)
							writer.Write([]byte("OK"))
						}).OnError(func(err error) {
							writer.WriteHeader(http.StatusInternalServerError)
							writer.Write([]byte("failed to write to kafka"))
							return
						}))
				} else {
					writer.WriteHeader(http.StatusBadRequest)
					writer.Write([]byte("provided id is not registered for provided type"))
				}
			}).OnError(func(err error) {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("failed to edit field: " + err.Error()))
			}))
		return behaviour
	}))
}

func (k *kanban) openStream(context context.Context, _ inspect.Inspectable, header http.Header, _ *multipart.Part, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")

	lastId := header["Last-Event-Id"]
	log.Println("got last id", lastId)
	startId := int64(0)
	if len(lastId) > 0 && len(lastId[0]) > 0 {
		id, err := strconv.ParseInt(lastId[0], 10, 64)
		if err == nil {
			startId = id + 1
		}
	}

	input := newStreamInput(writer)
	input.CloseWhenActorCloses()
	k.System().Become(shutdownactor.NewShutdownableActor(context.Done(), "http_stream",
		func(actor *actors.Actor) actors.Behaviour {
			onQuit := func(err error) {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(err.Error()))
				actor.Quit(err)
			}
			behaviour := actors.Behaviour{Name: "client stream handler"}
			actor.RequestStream(input, k.Service(), &subscribe{startId: startId}, onQuit)
			return behaviour
		}))
}

func init() {
	dir, _ := os.Getwd()
	defaultHttpServerParams := flags.HostPort{Port: 8882}
	defaultKafkaParams := flags.HostPort{Port: 9092}
	var keyValueTopic string
	var historyTopic string
	starter.SetCreator(KanbanName, func(s *actors.Actor, name string) (actors.ActorService, error) {
		return newKanban(KanbanName, s.System(), dir, defaultHttpServerParams, []string{defaultKafkaParams.String()}, keyValueTopic, historyTopic)
	})

	starter.SetFlagInitializer(KanbanName, func() {
		flags.StringFlag(&dir, "file-dir", "work directory")
		defaultHttpServerParams.RegisterFlagsWithDescriptions(
			"http",
			"listen to http requests on this hostname/ip address",
			"listen to http requests on this port")
		defaultKafkaParams.RegisterFlagsWithDescriptions(
			"kafka",
			"kafka hostname/ip address",
			"kafka port")
		flags.StringFlag(&keyValueTopic, "keyvalue", "kafka key value topic name")
		flags.StringFlag(&historyTopic, "history", "kafka history topic name")
	})
}
