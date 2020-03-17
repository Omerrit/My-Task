package auth

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/registry"
	"gerrit-share.lan/go/common"
)

const authServiceName = "auth"

type Handle struct {
	handle registry.Handle
}

func (h *Handle) Acquire(me actors.ActorCompatible, onFinished common.SimpleCallback, onError common.ErrorCallback) {
	h.handle.Acquire(authServiceName, me, onFinished, onError)
}

func (h *Handle) ConnectionEstablished(id Id) {
	h.handle.Send(&connectionEstablished{id})
}

func (h *Handle) ConnectionClosed(id Id) {
	h.handle.Send(&connectionClosed{id})
}

func (h *Handle) RequestUserId(connId Id, onFinished IdCallback, onError common.ErrorCallback) {
	h.handle.Request(&userRequest{connId},
		actors.NewReplyProcessor(func(reply interface{}) {
			onFinished.Call(reply.(*User).Id)
		}, onError))
}

func (h *Handle) RequestPermission(connId Id, command string, path string, onFinished func(bool), onError common.ErrorCallback) {
	h.handle.Request(&getPermission{connId, path, command}, actors.NewReplyProcessor(func(reply interface{}) {
		if onFinished == nil {
			return
		}
		onFinished(reply.(bool))
	}, onError))
}
