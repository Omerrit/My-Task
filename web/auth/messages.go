package auth

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

const packageName = "auth"

//notification message from auth service to auth filter
type connectionClosed struct {
	Id
}

//notification message from protocol service to auth service
type connectionEstablished struct {
	Id
}

//auth service->auth filter
type connectionPrivilegesChanged struct {
	Id
}

func init() {
	inspectables.RegisterDescribed(packageName+".newconn", func() inspect.Inspectable { return new(connectionEstablished) },
		"'connection established' message to be sent by protocol service to auth service")
	inspectables.RegisterDescribed(packageName+".connclosed", func() inspect.Inspectable { return new(connectionClosed) },
		"'connection closed' message, protocol service sends it to auth service and auth service resends it to auth filter clients that are waiting for it")
	inspectables.RegisterDescribed(packageName+".provchange", func() inspect.Inspectable { return new(connectionPrivilegesChanged) },
		"message indicates that privileges were changed for this connection just now, auth filter may re-request new privileges or just clear cached ones")
}
