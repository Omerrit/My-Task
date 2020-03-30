package registry

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/interfaces"
)

type command struct {
	cmd       inspect.Inspectable
	processor interfaces.ReplyProcessor
}

type Handle struct {
	Me            *actors.Actor
	Service       actors.ActorService
	messageBuffer []inspect.Inspectable
	commandBuffer []command
}

func (h *Handle) IsValid() bool {
	return h.Service != nil
}

func (h *Handle) Send(message inspect.Inspectable) {
	if h.Me == nil {
		return
	}
	if h.Service == nil {
		h.messageBuffer = append(h.messageBuffer, message)
		return
	}
	h.Me.SendMessage(h.Service, message)
}

func (h *Handle) Request(cmd inspect.Inspectable, processor interfaces.ReplyProcessor) {
	if h.Me == nil {
		processor.Error(actors.ErrApiHandleInvalid)
		return
	}
	if h.Service == nil {
		h.commandBuffer = append(h.commandBuffer, command{cmd, processor})
		return
	}
	h.Me.SendRequest(h.Service, cmd, processor)
}

func (h *Handle) Acquire(serviceName string, me actors.ActorCompatible, onFinished common.SimpleCallback, onError common.ErrorCallback) {
	h.Me = me.GetBase()
	WaitActor(me, serviceName, func(actor actors.ActorService) {
		h.Service = actor
		h.Me.Monitor(actor, func(error) {
			h.Service = nil
			h.Me = nil
		})
		onFinished.Call()
		if len(h.messageBuffer) > 0 {
			h.Me.SendMessages(h.Service, h.messageBuffer)
			h.messageBuffer = nil
		}
		for _, cmd := range h.commandBuffer {
			h.Me.SendRequest(h.Service, cmd.cmd, cmd.processor)
		}
		h.commandBuffer = nil
	}, func(err error) {
		h.Me = nil
		onError.Call(err)
		h.messageBuffer = nil
		h.commandBuffer = nil
	})
}
