package actors

import (
	"gerrit-share.lan/go/actors/internal/queue"
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/debug"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"gerrit-share.lan/go/interfaces"
)

type service struct {
	queue           *queue.Queue
	shutdownChannel common.SignalChannel
	err             error
	actor           ActorCompatible
}

const ActorServiceName = packageName + ".actor"

func zombieBehaviour(service ActorService, message interface{}, deadError error) {
	switch msg := message.(type) {
	case commandMessage:
		msg.replyWithError(deadError)
	case establishLink:
		if msg.linkType == linkMonitor {
			if service == nil {
				enqueue(msg.source, notifyClose{service, nil})
			} else {
				enqueue(msg.source, notifyClose{service, service.CloseError()})
			}
		} else {
			enqueue(msg.source, quitMessage{nil})
		}
	case streamCanSend:
		enqueue(msg.source, downstreamStopped{outputId{msg.streamId, service}, deadError})
	case streamRequest:
		enqueue(msg.id.destination, upstreamStopped{sourceId{msg.id.streamId, service}, deadError})
	case streamReply:
		enqueue(msg.id.source, downstreamStopped{outputId{msg.id.streamId, service}, deadError})
	case streamAck:
		enqueue(msg.id.destination, upstreamStopped{sourceId{msg.id.streamId, service}, deadError})
	case reply:
		if request, ok := msg.data.(RequestStream); ok {
			req := request.getStreamRequest()
			enqueue(req.id.destination, upstreamStopped{sourceId{req.id.streamId, service}, deadError})
		}
	case subscribeStateChange:
		enqueue(msg.source, notifyStateChange{service, ActorDead})
	}
}

func processMessagesAsZombie(service ActorService, head *queue.QueueElement, deadError error) {
	if head == nil {
		return
	}
	for ; head != nil; head = head.Next {
		if array, ok := head.Data.([]interface{}); ok {
			for _, item := range array {
				zombieBehaviour(service, item, deadError)
			}
		} else {
			zombieBehaviour(service, head.Data, deadError)
		}
	}
}

func (s *service) init(actor ActorCompatible, queue *queue.Queue) {
	s.actor = actor
	if s.queue == nil {
		s.queue = queue
		s.shutdownChannel = make(common.SignalChannel)
	}
}

func (s *service) enqueue(message interface{}) {
	if s == nil || s.queue == nil {
		if array, ok := message.([]interface{}); ok {
			for _, item := range array {
				zombieBehaviour(nil, item, ErrActorNull)
			}
		} else {
			zombieBehaviour(nil, message, ErrActorNull)
		}
		return
	}
	if !s.queue.Push(message) {
		if array, ok := message.([]interface{}); ok {
			for _, item := range array {
				zombieBehaviour(s.actor.GetBase().Service(), item, ErrActorDead)
			}
		} else {
			zombieBehaviour(s.actor.GetBase().Service(), message, ErrActorDead)
		}
		return
	}
	debug.Printf("+ %p: %#v\n", s.actor.GetBase().Service(), message)
}

func enqueue(service ActorService, message interface{}) {
	if service == nil {
		if array, ok := message.([]interface{}); ok {
			for _, item := range array {
				zombieBehaviour(nil, item, ErrActorNull)
			}
		} else {
			zombieBehaviour(nil, message, ErrActorNull)
		}
		return
	}
	service.enqueue(message)
}

func (s *service) SendMessage(message inspect.Inspectable) {
	s.enqueue(message)
}

func (s *service) DoneChannel() common.OutSignalChannel {
	return s.shutdownChannel.ToOutput()
}

func (s *service) Shutdown() {
	s.enqueue(closeMessage{nil})
}

func (s *service) Close() error {
	s.Shutdown()
	<-s.shutdownChannel
	return s.err
}

func (s *service) SendQuit(err error) {
	s.enqueue(quitMessage{err})
}

func (s *service) CloseError() error {
	return s.err
}

func (s *service) System() *System {
	return s.actor.GetBase().System()
}

func (s *service) start(system *System) {
	var result error
	defer func() {
		err := errors.MakeArray(result, errors.RecoverToError(recover()))
		base := s.actor.GetBase()
		base.close()
		if shutdowner, ok := s.actor.(interfaces.Shutdownable); ok {
			err.Add(shutdowner.Shutdown())
		}
		s.err = err.ToError()
		processMessagesAsZombie(base.Service(), s.queue.TakeHeadAndClose(), ErrActorDead)
		close(s.shutdownChannel)
		debug.Printf("%p died\n", base.Service())
		system.serviceFinished()
	}()
	if behavioral, ok := s.actor.(BehavioralActor); ok {
		behavioral.GetBase().setBehaviour(system, behavioral.MakeBehaviour())
	}
	s.enqueue(initialMessage{})
	result = s.actor.Run()
	if result == nil {
		result = s.actor.GetBase().quitError
	}
}

func (s *service) getActor() ActorCompatible {
	return s.actor
}

func (s *service) Visit(visitor ResponseVisitor) {
	visitor.Reply(s)
}

//don't call this directly, use InspectActorService when necessary
func (s *service) Inspect(inspector *inspect.GenericInspector) {
	//STUB: do actual implementation
	inspector.Object(ActorServiceName, "").End()
}

//should be called to inspect actor service variable, will initialize it if it's nil
func InspectActorService(actor *ActorService, inspector *inspect.GenericInspector) {
	if inspector.IsReading() {
		if actor == nil {
			*actor = new(service)
		}
	}
	(*actor).Inspect(inspector)
}

func init() {
	inspectables.Register(ActorServiceName, func() inspect.Inspectable { return new(service) })
}
