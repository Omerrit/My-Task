package example

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/registry"
	"gerrit-share.lan/go/common"
)

const serviceName = "example_service"

type Handle struct {
	handle registry.Handle
}

//mandatory initialization function
func (h *Handle) Acquire(me actors.ActorCompatible, onFinished common.SimpleCallback, onError common.ErrorCallback) {
	h.handle.Acquire(serviceName, me, onFinished, onError)
}

//always implement this function
func (h *Handle) IsValid() bool {
	return h.handle.IsValid()
}

func (h *Handle) GetStuff(me actors.ActorCompatible, path string, onFinished func(string, int), onError common.ErrorCallback) {
	//use registry.Handle functions Request and Send to make requests and send commands, it buffers them if service is still waiting for an actor
	//Do not abuse it though, buffer could become very large if actor never starts
	//it's safer to wait for initialization to finish before making requests
	if onFinished == nil {
		h.handle.Request(&getStuff{pathHolder{path}}, actors.NewErrorProcessor(onError))
		return
	}
	h.handle.Request(&getStuff{pathHolder{path}}, actors.NewReplyProcessor(func(reply interface{}) {
		value := reply.(*stuff)
		onFinished(value.name, value.value)
	}, onError))
}
