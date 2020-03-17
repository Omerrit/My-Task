package starter

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/registry"
	"gerrit-share.lan/go/common"
)

const starterServiceName = "starter"

type Handle struct {
	handle registry.Handle
}

func (h *Handle) Acquire(me actors.ActorCompatible, onFinished common.SimpleCallback, onError common.ErrorCallback) {
	h.handle.Acquire(starterServiceName, me, onFinished, onError)
}

func (h *Handle) IsValid() bool {
	return h.handle.IsValid()
}

func (h *Handle) DependOn() {
	if h.handle.IsValid() {
		h.handle.Me.DependOn(h.handle.Service)
	}
}

func (h *Handle) Link() {
	if h.handle.IsValid() {
		h.handle.Me.Link(h.handle.Service)
	}
}
