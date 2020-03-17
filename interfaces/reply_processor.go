package interfaces

import ()

type ErrorProcessor interface {
	Error(error)
}

type ReplyProcessor interface {
	Process(interface{})
	ErrorProcessor
}

type DummyReplyProcessor struct{}

func (DummyReplyProcessor) Process(interface{}) {}
func (DummyReplyProcessor) Error(error)         {}

func TestReplyProcessor(ReplyProcessor) {}
