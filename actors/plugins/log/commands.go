package log

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

type subscribe struct {
	actors.RequestStreamBase
	id int
}

const subscribeName = packageName + ".subscribe"

func (s *subscribe) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(subscribeName, "")
	{
		s.RequestStreamBase.Inspect(objectInspector.Value("requeststream", true, ""))
		objectInspector.Int(&s.id, "id", true, "")
		objectInspector.End()
	}
}

func init() {
	inspectables.RegisterDescribed(subscribeName, func() inspect.Inspectable { return new(subscribe) }, "subscribe to logger message stream")
}

type notGonnaSubscribe struct {
	id int
}

const notGonnaSubscribeName = packageName + ".nosubscribe"

func (n *notGonnaSubscribe) Inspect(i *inspect.GenericInspector) {
	i.Int(&n.id, notGonnaSubscribeName, "")
}

func init() {
	inspectables.RegisterDescribed(notGonnaSubscribeName, func() inspect.Inspectable { return new(notGonnaSubscribe) }, "notify logger writer isn't subscribing")
}
