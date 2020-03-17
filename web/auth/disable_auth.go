package auth

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/registry"
)

func DisableAuth(me actors.ActorCompatible) {
	registry.MarkAsUnregisterable(me, authServiceName, nil)
}
