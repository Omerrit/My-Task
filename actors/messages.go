package actors

import (
	"gerrit-share.lan/go/inspect"
)

type initialMessage struct {
}

type closeMessage struct {
	err error
}

type quitMessage struct {
	err error
}

type notifyClose struct {
	destination ActorService
	err         error
}

type establishLink struct {
	source   ActorService
	linkType linkType
}

type commandId int64

type promiseId struct {
	origin ActorService
	id     commandId
}

type commandMessage struct {
	promiseId
	data inspect.Inspectable
}

//i'm processing this command!
type preReply struct {
	id        commandId
	processor ActorService
}

type cancelCommand struct {
	origin ActorService
	id     commandId
}

type reply struct {
	id   commandId
	data interface{}
}

type errorReply struct {
	err error
}

func (p promiseId) makeReply(data interface{}) reply {
	return reply{p.id, data}
}

func (p promiseId) makeErrorReply(err error) reply {
	return p.makeReply(errorReply{err})
}

func (p promiseId) reply(data interface{}) {
	if p.origin != nil {
		enqueue(p.origin, p.makeReply(data))
	}
}

func (p promiseId) replyWithError(err error) {
	if p.origin != nil {
		enqueue(p.origin, p.makeErrorReply(err))
	}
}

func (c commandMessage) preReply(me ActorService) {
	if c.origin != nil {
		enqueue(c.origin, preReply{c.id, me})
	}
}

func (c *commandMessage) invalidate() {
	c.origin = nil
}

func (c *commandMessage) isValid() bool {
	return c.origin != nil
}

func (c cancelCommand) toPromiseId() promiseId {
	return promiseId{c.origin, c.id}
}
