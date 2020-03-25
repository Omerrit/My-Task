package actortest

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/errors"
	"testing"
)

func EnsureDead(t *testing.T, service actors.ActorService, expectedErrors ...error) {
	service.System().Become(actors.NewSimpleActor(func(actor *actors.Actor) actors.Behaviour {
		actor.SetFinishedServiceProcessor(func(dead actors.ActorService, err error) {
			if dead != service {
				t.Error("got actor dead message for the actor I didn't monitor")
				actor.Quit(nil)
				return
			}
			for _, e := range expectedErrors {
				if !errors.Is(err, e) {
					t.Error("actor closed with unexpected error:", err)
					var ste errors.StackTraceError
					if errors.As(err, &ste) {
						t.Error(ste.StackTrace())
					}
				}
			}
			//actor should quit by itself
		})
		actor.Monitor(service)
		return actors.Behaviour{}
	}))
}
