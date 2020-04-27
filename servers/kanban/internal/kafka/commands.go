package kafka

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

const packageName = "kafka"

type Subscribe struct {
	actors.RequestStreamBase
}

const SubscribeName = packageName + ".subscribe"

func init() {
	inspectables.RegisterDescribed(SubscribeName, func() inspect.Inspectable { return new(Subscribe) }, "subscribe to kafka stream")
}
