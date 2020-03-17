package published

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

const packageName = "actplugins.published"

type publish struct {
	actor actors.ActorService
}

func (p *publish) Inspect(inspector *inspect.GenericInspector) {
	actors.InspectActorService(&p.actor, inspector)
}

type subscribe struct {
	actors.RequestStreamBase
}

func init() {
	inspectables.RegisterDescribed(packageName+".publish", func() inspect.Inspectable { return new(publish) }, "publish actor")
	inspectables.RegisterDescribed(packageName+".subscribe", func() inspect.Inspectable { return new(subscribe) }, "subscribe to published actor stream")
}
