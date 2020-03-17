package registry

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
		RegisterMe(actor, testName, actortest.QuitOnError(t, actor, "registration failed"))
		GetActor(actor, testName, func(service actors.ActorService) {
			if service != actor.Service() {
				t.Error("got wrong service")
			}
		}, actortest.QuitOnError(t, actor, "failed to get actor"))
		return actors.Behaviour{}
	}))
	actortest.EnsureDead(t, me, nil)
}

func TestWait(t *testing.T) {
	var s actors.System
	defer s.WaitFinished()
	me := s.Become(actortest.NewTestingActor(t, func(actor *actors.Actor) actors.Behaviour {
		WaitActor(actor, testName, func(service actors.ActorService) {
			if service != actor.Service() {
				t.Error("got wrong service")
			}
		}, actortest.QuitOnError(t, actor, "failed to get actor"))
		RegisterMe(actor, testName, actortest.QuitOnError(t, actor, "registration failed"))
		return actors.Behaviour{}
	}))
	actortest.EnsureDead(t, me, nil)
}

func TestSubscribe(t *testing.T) {
	var s actors.System
	defer s.WaitFinished()
	me := s.Become(actortest.NewTestingActor(t, func(actor *actors.Actor) actors.Behaviour {
		var in *actors.SimpleCallbackStreamInput
		in = actors.NewSimpleCallbackStreamInput(func(data inspect.Inspectable) error {
			info := data.(*Info)
			if info.Name != testName {
				t.Error("got wrong name:", info.Name, ",expected:", testName)
			}
			if info.Actor != actor.Service() {
				t.Error("got wrong service")
			}
			in.Close()
			return nil
		}, func(base *actors.StreamInputBase) {
			base.RequestData(new(Info), 1)
		})
		Subscribe(actor, in, actortest.QuitOnError(t, actor, "subscription failed"))
		RegisterMe(actor, testName, actortest.QuitOnError(t, actor, "registration failed"))
		return actors.Behaviour{}
	}))
	actortest.EnsureDead(t, me, nil)
}
