package shutdownactor

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/common"
)

type shutdownableActor struct {
	actors.Actor
	makeBehaviour   func(*actors.Actor) actors.Behaviour
	shutdownChannel common.OutSignalChannel
	name            string
}

func (s *shutdownableActor) Run() error {
	for {
		select {
		case <-s.shutdownChannel:
			s.Quit(nil)
			if !s.ProcessMessages() {
				return nil
			}
			return s.Actor.Run()
		case <-s.IncomingChannel():
			if !s.ProcessMessages() {
				return nil
			}
		}
	}
}

func (s *shutdownableActor) MakeBehaviour() actors.Behaviour {
	if s.makeBehaviour != nil {
		b := s.makeBehaviour(&s.Actor)
		b.Name = s.name
		return b
	}
	return actors.Behaviour{Name: s.name}
}

func NewShutdownableActor(shutdown common.OutSignalChannel, name string, behaviourMaker actors.BehaviourMaker) actors.BehavioralActor {
	if behaviourMaker == nil {
		return nil
	}
	return &shutdownableActor{
		shutdownChannel: shutdown,
		makeBehaviour:   behaviourMaker,
		name:            name,
	}
}
