package actors

import (
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/utils/callbackarrays"
)

type StreamInput interface {
	getBase() *StreamInputBase
	Process(data inspect.Inspectable) error
	RequestNext()
}

type StreamInputBase struct {
	id                         streamId
	source                     ActorService
	me                         *Actor
	isSuspended                bool
	isWaitingForData           bool
	shouldClose                bool
	shouldCloseWhenActorCloses bool
	request                    streamRequest
	onClose                    callbackarrays.ErrorCallbacks
}

func (s *StreamInputBase) getBase() *StreamInputBase {
	return s
}

func (s *StreamInputBase) init(me *Actor, id streamId) {
	s.id = id
	s.me = me
	s.request.id.streamId = s.id
	s.request.id.destination = s.me.Service()
}

func (s *StreamInputBase) setSource(source ActorService) {
	s.source = source
	if s.shouldClose {
		s.Close()
		s.shouldClose = false
	}
}

func (s *StreamInputBase) OnClose(callback common.ErrorCallback) {
	s.onClose.Push(callback)
}

func (s *StreamInputBase) closed(err error) {
	s.onClose.Run(err)
	s.id.Invalidate()
}

func (s *StreamInputBase) Close() {
	if s.source == nil {
		s.shouldClose = true
		return
	}
	enqueue(s.source, closeStream{s.id, s.me.Service()})
}

func (s *StreamInputBase) CloseWhenActorCloses() {
	s.shouldCloseWhenActorCloses = true
	if s.me != nil && s.me.state != ActorRunning {
		s.Close()
	}
}

func (s *StreamInputBase) Suspend() {
	s.isSuspended = true
}

func (s *StreamInputBase) Resume() {
	s.isSuspended = false
	if s.source != nil && !s.isWaitingForData && s.request.data != nil {
		enqueue(s.source, s.request)
		s.request.data = nil
		s.isWaitingForData = true
	}
}

func (s *StreamInputBase) IsSuspended() bool {
	return s.isSuspended
}

func (s *StreamInputBase) RequestData(data inspect.Inspectable, maxLen int) {
	if !s.isWaitingForData {
		if s.isSuspended || s.source == nil {
			s.request.data = data
			s.request.maxLen = maxLen
		} else {
			enqueue(s.source, streamRequest{outputId{s.id, s.me.Service()}, data, maxLen})
			s.request.data = nil
			s.isWaitingForData = true
		}
	}
}

func (s *StreamInputBase) Acknowledge() {
	if !s.isWaitingForData && !s.isSuspended && s.source != nil {
		enqueue(s.source, streamAck{outputId{s.id, s.me.Service()}})
	}
}

func (s *StreamInputBase) dataReceived() {
	s.isWaitingForData = false
}
