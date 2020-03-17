package actorrp

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/interfaces"
)

type simpleActorReplyProcessor actors.ActorServiceCallback

func (s simpleActorReplyProcessor) Process(reply interface{}) {
	if reply == nil {
		s(nil)
		return
	}
	s(reply.(actors.ActorService))
}

func (s simpleActorReplyProcessor) Error(err error) {}

type actorReplyProcessor struct {
	simpleActorReplyProcessor
	onError common.ErrorCallback
}

func (a *actorReplyProcessor) Error(err error) {
	a.onError(err)
}

func New(onReply actors.ActorServiceCallback, onError common.ErrorCallback) interfaces.ReplyProcessor {
	if onReply == nil {
		return actors.NewErrorProcessor(onError)
	}
	if onError == nil {
		return simpleActorReplyProcessor(onReply)
	}
	return &actorReplyProcessor{simpleActorReplyProcessor(onReply), onError}
}
