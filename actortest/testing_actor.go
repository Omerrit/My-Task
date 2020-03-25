package actortest

import (
	"gerrit-share.lan/go/actors"
	"testing"
)

func NewTestingActor(t *testing.T, behaviourMaker actors.BehaviourMaker) actors.BehavioralActor {
	return actors.NewSimpleActor(func(actor *actors.Actor) actors.Behaviour {
		PrintOnPanic(t, actor)
		if behaviourMaker != nil {
			return behaviourMaker(actor)
		}
		return actors.Behaviour{}
	})
}
