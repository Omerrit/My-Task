package actors

import (
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/interfaces"
)

type replyCallbackProcessor struct {
	onProcess common.GenericCallback
	onError   common.ErrorCallback
}

func (r replyCallbackProcessor) Process(data interface{}) {
	r.onProcess(data)
}

func (r replyCallbackProcessor) Error(err error) {
	r.onError(err)
}

type simpleReplyCallbackProcessor common.GenericCallback

func (s simpleReplyCallbackProcessor) Process(data interface{}) {
	s(data)
}

func (s simpleReplyCallbackProcessor) Error(err error) {
}

func (s simpleReplyCallbackProcessor) OnError(onError common.ErrorCallback) *replyCallbackProcessor {
	return &replyCallbackProcessor{common.GenericCallback(s), onError}
}

func (s simpleReplyCallbackProcessor) QuitOnError(me ActorCompatible, description string) *replyCallbackProcessor {
	return &replyCallbackProcessor{common.GenericCallback(s), func(err error) {
		me.GetBase().Quit(errors.Describe(err, description))
	}}
}

func (s simpleReplyCallbackProcessor) QuitOnErrorf(me ActorCompatible, format string, args ...interface{}) *replyCallbackProcessor {
	return &replyCallbackProcessor{common.GenericCallback(s), func(err error) {
		me.GetBase().Quit(errors.Describef(err, format, args...))
	}}
}

func OnReply(onProcess common.GenericCallback) simpleReplyCallbackProcessor {
	return simpleReplyCallbackProcessor(onProcess)
}

type errorReplyCallbackProcessor common.ErrorCallback

func (e errorReplyCallbackProcessor) Process(data interface{}) {
}

func (e errorReplyCallbackProcessor) Error(err error) {
	e(err)
}

func (e errorReplyCallbackProcessor) OnReply(onProcess common.GenericCallback) *replyCallbackProcessor {
	return &replyCallbackProcessor{onProcess, common.ErrorCallback(e)}
}

func OnReplyError(onError common.ErrorCallback) errorReplyCallbackProcessor {
	return errorReplyCallbackProcessor(onError)
}

func OnReplyErrorQuit(me ActorCompatible, description string) errorReplyCallbackProcessor {
	return func(err error) {
		me.GetBase().Quit(errors.Describe(err, description))
	}
}

func OnReplyErrorQuitf(me ActorCompatible, format string, args ...interface{}) errorReplyCallbackProcessor {
	return func(err error) {
		me.GetBase().Quit(errors.Describef(err, format, args...))
	}
}

//chooses which real reply processor implementation to use based on nil callbacks.
//If both callbacks are nil returns nil to signal that actor shouldn't wait for reply.
//Command cancellation requires request to be tracked anyway
func NewReplyProcessor(onProcess common.GenericCallback, onError common.ErrorCallback) interfaces.ReplyProcessor {
	if onProcess == nil {
		if onError == nil {
			return nil
		}
		return OnReplyError(onError)
	}
	if onError == nil {
		return OnReply(onProcess)
	}
	return &replyCallbackProcessor{onProcess, onError}
}

//same as NewReplyProcessor but without the ability to set onProcess at all
func NewErrorProcessor(onError common.ErrorCallback) interfaces.ReplyProcessor {
	if onError == nil {
		return nil
	}
	return errorReplyCallbackProcessor(onError)
}
