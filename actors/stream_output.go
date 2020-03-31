package actors

import (
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/utils/callbackarrays"
)

type StreamOutput interface {
	getBase() *StreamOutputBase
	DataSource
}

func TestStreamOutput(StreamOutput) {}

type StreamOutputSet map[StreamOutput]common.None

func (s *StreamOutputSet) Add(out StreamOutput) {
	if *s == nil {
		*s = make(StreamOutputSet, 1)
	}
	(*s)[out] = common.None{}
}

func (s *StreamOutputSet) Remove(out StreamOutput) {
	delete(*s, out)
}

func (s *StreamOutputSet) Contains(out StreamOutput) bool {
	_, ok := (*s)[out]
	return ok
}

func (s *StreamOutputSet) IsEmpty() bool {
	return len(*s) == 0
}

func (s *StreamOutputSet) Clear() {
	*s = nil
}

type StreamOutputBase struct {
	outStreamId                 outputId
	actor                       *Actor
	onClose                     callbackarrays.ErrorCallbacks
	isStreamClosing             bool
	closeError                  error
	isFlushScheduled            bool
	shouldCloseWhenAcknowledged bool
	shouldCloseWhenActorCloses  bool
}

func (d *StreamOutputBase) getBase() *StreamOutputBase {
	return d
}

func (d *StreamOutputBase) init(actor *Actor, id outputId) {
	d.actor = actor
	d.outStreamId = id
}

func (d *StreamOutputBase) FlushLater() {
	if d.isFlushScheduled || d.actor == nil {
		return
	}
	d.actor.markOutputReady(d.outStreamId)
	d.isFlushScheduled = true
}

func (d *StreamOutputBase) OnClose(processor common.ErrorCallback) {
	d.onClose.PushNonNull(processor)
}

func (d *StreamOutputBase) acknowledged() {
	d.isFlushScheduled = false
	if d.shouldCloseWhenAcknowledged {
		d.CloseStream(d.closeError)
		d.shouldCloseWhenAcknowledged = true
	}
}

func (d *StreamOutputBase) CloseStream(err error) {
	if d.isStreamClosing {
		return
	}
	d.isStreamClosing = true
	d.FlushLater()
	d.closeError = err
}

func (d *StreamOutputBase) closeStreamNow(err error) {
	d.isStreamClosing = true
	d.closeError = err
}

func (d *StreamOutputBase) streamClosed(err error) {
	d.actor = nil
	d.onClose.Run(err)
}

func (d *StreamOutputBase) IsValid() bool {
	return d.actor != nil
}

func (d *StreamOutputBase) CloseWhenAcknowledged() {
	d.shouldCloseWhenAcknowledged = true
}

func (d *StreamOutputBase) CloseWhenActorCloses() {
	d.shouldCloseWhenActorCloses = true
	if d.actor != nil && d.actor.state != ActorRunning {
		d.CloseStream(d.actor.quitError)
	}
}

type genericCallbackStreamOutput struct {
	StreamOutputBase
	generator    func(*StreamOutputBase, inspect.Inspectable, int) (inspect.Inspectable, error)
	acknowledged common.SimpleCallback
}

func (g *genericCallbackStreamOutput) FillData(data inspect.Inspectable, maxLen int) (inspect.Inspectable, error) {
	if g.generator != nil {
		return g.generator(&g.StreamOutputBase, data, maxLen)
	}
	return nil, nil
}

func (g *genericCallbackStreamOutput) Acknowledged() {
	if g.acknowledged != nil {
		g.acknowledged()
	}
}
