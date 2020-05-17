package kanban

import (
	"context"
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/replies"
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
	"github.com/Shopify/sarama"
	"io/ioutil"
	"log"
	"net/http"
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
)

type kanban struct {
	actors.Actor
	server       http.Server
	name         string
	httpActor    actors.ActorService
	broadcaster  actors.StateBroadcaster
	messages     kafka.MessagesStream
	config       *sarama.Config
	brokers      []string
	topic        string
	producer     sarama.SyncProducer
	usersStorage utils.UsersStorage
	endpoints    endpoints.Endpoints
	ids          ids.TypedIds
	topicKeys    sets.String
	configLength int
	eofReached   bool
}

func newKanban(name string, system *actors.System, httpHostPort flags.HostPort, kafkaBrokers []string, topic string) (actors.ActorService, error) {
	server := new(kanban)
	server.name = name
	server.topic = topic
	server.configLength = len(typeMessages) + len(propMessages)
	server.brokers = kafkaBrokers
	server.config = kafka.NewConfig()
	server.generateEndpoints()
	server.server = http.Server{
		Addr:    httpHostPort.String(),
		Handler: server}
	err := kafka.CheckTopic(server.topic, server.brokers, server.config, kafka.CompactTopicEntries())
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
		if !k.eofReached {
			return nil, fmt.Errorf("restoring from kafka is in progress")
		}
		newIdCmd := cmd.(*newId)
		return replies.String(k.ids.AcquireNewId(newIdCmd.objectType, newIdCmd.id)), nil
	}).ResultString()
	behaviour.AddCommand(new(deleteId), func(cmd interface{}) (actors.Response, error) {
		if !k.eofReached {
			return nil, fmt.Errorf("restoring from kafka is in progress")
		}
		deleteCmd := cmd.(*deleteId)
		return nil, k.ids.DeleteId(deleteCmd.objectType, deleteCmd.id)
	})
	behaviour.AddCommand(new(isIdRegistered), func(cmd interface{}) (actors.Response, error) {
		if !k.eofReached {
			return nil, fmt.Errorf("restoring from kafka is in progress")
		}
		isRegisteredCmd := cmd.(*isIdRegistered)
		return replies.Bool(k.ids.IsRegistered(isRegisteredCmd.objectType, isRegisteredCmd.id)), nil
	}).ResultBool()
	behaviour.AddCommand(new(reserveId), func(cmd interface{}) (actors.Response, error) {
		if !k.eofReached {
			return nil, fmt.Errorf("restoring from kafka is in progress")
		}
		reserveIdCmd := cmd.(*reserveId)
		return nil, k.ids.RestoreId(reserveIdCmd.objectType, reserveIdCmd.id)
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
		return &message{}
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
}

func (k *kanban) consumeFromKafka(offset int64, processor func(*kafka.Message) (error, bool)) error {
	consumer, err := kafka.NewConsumer("kafka_consumer", k.System(), k.config, k.brokers, k.topic, 0, offset)
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

func (k *kanban) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	if request.Method == "OPTIONS" {
		writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		writer.WriteHeader(http.StatusOK)
		return
	}

	endpoint, ok := k.endpoints[strings.TrimSuffix(request.URL.Path, "/")]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	switch contentType := request.Header.Get("Content-Type"); contentType {
	case "application/json":
		k.processJsonRequest(request, endpoint, writer)
	default:
		if endpoint.Creator() != nil {
			writer.WriteHeader(http.StatusUnsupportedMediaType)
			writer.Write([]byte("unsupported content type"))
			return
		}
		endpoint.Handler()(nil, request.Header, writer)
	}
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
	endpoint.Handler()(command, request.Header, writer)
}

func (k *kanban) reserveId(command inspect.Inspectable, _ http.Header, writer http.ResponseWriter) {
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

func (k *kanban) deleteId(command inspect.Inspectable, _ http.Header, writer http.ResponseWriter) {
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

func (k *kanban) newId(command inspect.Inspectable, _ http.Header, writer http.ResponseWriter) {
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

func (k *kanban) login(command inspect.Inspectable, _ http.Header, writer http.ResponseWriter) {
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

func (k *kanban) sendRestOfConfig() error {
	sendMessage := func(key, value string) error {
		_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
			Topic:     k.topic,
			Key:       sarama.StringEncoder(key),
			Value:     sarama.StringEncoder(value),
			Headers:   nil,
			Metadata:  nil,
			Offset:    0,
			Partition: 0,
			Timestamp: time.Now(),
		})
		return err
	}

	for key, value := range typeMessages {
		err := sendMessage(key, value)
		if err != nil {
			return err
		}
	}
	for key, value := range propMessages {
		err := sendMessage(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *kanban) writeMessage(command *kafka.Message) (error, bool) {
	if command.Offset == kafka.OffsetUninitialized {
		k.eofReached = true
		if !k.isTopicValid() {
			err := fmt.Errorf("shitty topic detected! use empty or properly configured topic next time")
			log.Println(err)
			k.Quit(err)
			return err, false
		}
		k.topicKeys.Clear()
		k.runHttp()
		return k.sendRestOfConfig(), false
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

func (k *kanban) edit(command inspect.Inspectable, _ http.Header, writer http.ResponseWriter) {
	msgCommand := command.(*message)
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
					var value sarama.Encoder
					if len(msgCommand.value) > 0 {
						value = sarama.StringEncoder(msgCommand.value)
					}

					_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
						Topic:     k.topic,
						Key:       sarama.StringEncoder(msgCommand.key),
						Value:     value,
						Headers:   nil,
						Metadata:  nil,
						Offset:    0,
						Partition: 0,
						Timestamp: time.Now(),
					})
					if err != nil {
						writer.WriteHeader(http.StatusInternalServerError)
						writer.Write([]byte("failed to write to kafka"))
						return
					}
					writer.WriteHeader(http.StatusOK)
					writer.Write([]byte("OK"))
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

func (k *kanban) openStream(_ inspect.Inspectable, header http.Header, writer http.ResponseWriter) {
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

	k.System().BecomeFunc(func(actor *actors.Actor) actors.Behaviour {
		onQuit := func(err error) {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			actor.Quit(err)
		}

		behaviour := actors.Behaviour{Name: "client stream handler"}
		actor.RequestStream(newStreamInput(writer), k.Service(), &subscribe{startId: startId}, onQuit)
		return behaviour
	})
}

func init() {
	defaultHttpServerParams := flags.HostPort{Port: 8882}
	defaultKafkaParams := flags.HostPort{Port: 9092}
	var topic string
	starter.SetCreator(KanbanName, func(s *actors.Actor, name string) (actors.ActorService, error) {
		return newKanban(KanbanName, s.System(), defaultHttpServerParams, []string{defaultKafkaParams.String()}, topic)
	})

	starter.SetFlagInitializer(KanbanName, func() {
		defaultHttpServerParams.RegisterFlagsWithDescriptions(
			"http",
			"listen to http requests on this hostname/ip address",
			"listen to http requests on this port")
		defaultKafkaParams.RegisterFlagsWithDescriptions(
			"kafka",
			"kafka hostname/ip address",
			"kafka port")
		flags.StringFlag(&topic, "topic", "kafka topic name")
	})
}
