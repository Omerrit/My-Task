package published

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actortest"
	"gerrit-share.lan/go/inspect"
	"testing"
)

const testName = "testactor"

func TestGet(t *testing.T) {
	var s actors.System
	defer s.WaitFinished()
	me := s.Become(actortest.NewTestingActor(t, func(actor *actors.Actor) actors.Behaviour {
		var in *actors.SimpleCallbackStreamInput
		in = actors.NewSimpleCallbackStreamInput(func(data inspect.Inspectable) error {
			array := *(data.(*actors.ActorsArray))
			if len(array) != 1 {
				t.Error("got wrong array length")
				in.Close()
				return nil
			}
			if array[0] != actor.Service() {
				t.Error("got wrong service")
			}
			in.Close()
			return nil
		}, func(base *actors.StreamInputBase) {
			base.RequestData(new(actors.ActorsArray), 1)
		})
		Subscribe(actor, in, actortest.QuitOnError(t, actor, "subscription failed"))
		Publish(actor, actortest.QuitOnError(t, actor, "registration failed"))
		return actors.Behaviour{}
	}))
	actortest.EnsureDead(t, me, nil)
}
