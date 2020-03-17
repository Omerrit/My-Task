package actorwrappers

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/common"
)

type shutdownableActor struct {
	actors.Actor
	makeBehaviour   func(*actors.Actor) actors.Behaviour
	shutdownChannel common.OutSignalChannel
	onShutdown      func()
}

func (s *shutdownableActor) Run() error {
	for {
		select {
		case <-s.shutdownChannel:
			s.onShutdown()
			return nil
		case <-s.IncomingChannel():
			if !s.ProcessMessages() {
				return nil
			}
		}
	}
}

func (s *shutdownableActor) MakeBehaviour() actors.Behaviour {
	if s.makeBehaviour != nil {
		return s.makeBehaviour(&s.Actor)
	}
	return actors.Behaviour{}
}

func NewShutdownableActor(shutdown common.OutSignalChannel, onShutdown func(), behaviourMaker actors.BehaviourMaker) actors.BehavioralActor {
	return &shutdownableActor{
		shutdownChannel: shutdown,
		onShutdown:      onShutdown,
		makeBehaviour:   behaviourMaker,
	}
}
