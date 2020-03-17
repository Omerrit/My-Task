package autorestarter

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

const packageName = "actors.autorestarter"

type autostartMessage struct {
	inspect.EmptyObject
}

func init() {
	inspectables.RegisterDescribed(packageName+".autostart", func() inspect.Inspectable { return autostartMessage{} }, "to be sent to self, triggers autorestart")
}
