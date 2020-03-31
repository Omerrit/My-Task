package actortest

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/errors"
	"testing"
)

func EnsureDead(t *testing.T, service actors.ActorService, expectedErrors ...error) {
	service.System().Become(actors.NewSimpleActor(func(actor *actors.Actor) actors.Behaviour {
		actor.Monitor(service, func(err error) {
			for _, e := range expectedErrors {
				if !errors.Is(err, e) {
					t.Error("actor closed with unexpected error:", errors.FullInfo(err))
				}
			}
			//actor should quit by itself
		})
		return actors.Behaviour{}
	}))
}
