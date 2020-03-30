package actors

import (
	"fmt"
	"gerrit-share.lan/go/debug"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
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
	debug.Printf("actor %p died", service)
	if err != nil {
		debug.Println(err)
	} else {
		debug.Println("")
	}
	var sterr errors.StackTraceError
	if errors.As(err, &sterr) {
		debug.Println(sterr.StackTrace())
	}
}

func (s *surveillanceActor) addActor(service ActorService) {
	debug.Printf("actor %p started\n", service)
	s.actors.Add(service)
	s.Monitor(service, func(err error) {
		s.onActorStopped(service, err)
	})
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
	return behaviour
}
