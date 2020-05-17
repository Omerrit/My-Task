package actors

import (
	"fmt"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"log"
)

type closeMe struct {
	DummyFiller
}

const closeMeName = packageName + ".closeme"

func (closeMe) Inspect(inspector *inspect.GenericInspector) {
	inspector.Object(closeMeName, "free request to close actor (actor dependent processing)").End()
}

func init() {
	inspectables.Register(closeMeName, func() inspect.Inspectable { return closeMe{} })
}

type requestActorStream struct {
	RequestStreamBase
}

func init() {
	inspectables.RegisterDescribed(packageName+".subscribe", func() inspect.Inspectable { return new(requestActorStream) },
		"Start actor stream to the requestor")
}

type surveillanceActor struct {
	Actor
	broadcaster StateBroadcaster
	actors      ActorStateChangeStream
	shouldClose bool
}

func (s *surveillanceActor) onActorStopped(service ActorService, err error) {
	s.actors.Remove(service)
	if s.shouldClose {
		s.tryClose()
	}
	log.Printf("actor %p[%s] died", service, service.Name())
	if err != nil {
		log.Println(err)
	} else {
		log.Println("")
	}
	var sterr errors.StackTraceError
	if errors.As(err, &sterr) {
		log.Println(sterr.StackTrace())
	}
}

func (s *surveillanceActor) addActor(service ActorService) {
	log.Printf("actor %p[%s] started\n", service, service.Name())
	s.actors.Add(service)
	s.Monitor(service, func(err error) {
		s.onActorStopped(service, err)
	})
	s.MonitorStateChanges(service)
	s.broadcaster.NewDataAvailable()
}

func (s *surveillanceActor) tryClose() {
	if s.actors.state.IsEmpty() {
		s.Quit(nil)
	}
}

func (s *surveillanceActor) addSubscriber(request *requestActorStream) {
	//output would be closed as it should if filling initial data fails during initialization
	s.InitStreamOutput(s.broadcaster.AddOutput(), request)
}

var stateMessages = map[ActorState]string{
	ActorRunning:  "running",
	ActorQuitting: "quitting",
	ActorClosed:   "closed",
	ActorDead:     "dead"}

func (s *surveillanceActor) stateChangeMonitor(service ActorService, state ActorState) {
	log.Printf("%p [%s]: %s\n", service, service.Name(), stateMessages[state])
}

func (s *surveillanceActor) MakeBehaviour() Behaviour {
	var behaviour Behaviour
	s.broadcaster = NewBroadcaster(&s.actors)
	s.broadcaster.CloseWhenActorCloses()
	behaviour.AddMessage(new(service), func(cmd interface{}) {
		s.addActor(cmd.(*service))
	}).AddMessage(closeMe{}, func(interface{}) {
		s.shouldClose = true
		s.tryClose()
	})

	behaviour.AddCommand(new(requestActorStream), func(cmd interface{}) (Response, error) {
		s.addSubscriber(cmd.(*requestActorStream))
		return nil, nil
	}).Result(new(ActorsArray))
	s.SetPanicProcessor(func(err errors.StackTraceError) {
		fmt.Println(err)
		fmt.Println(err.StackTrace())
		s.Quit(err)
	})
	s.SetStateChangeProcessor(s.stateChangeMonitor)
	behaviour.Name = "surveillance actor"
	return behaviour
}
