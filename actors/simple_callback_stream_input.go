package actors

import (
	"gerrit-share.lan/go/inspect"
)

type SimpleCallbackStreamInput struct {
	StreamInputBase
	processor func(inspect.Inspectable) error
	requestor func(*StreamInputBase)
}

func NewSimpleCallbackStreamInput(processor func(inspect.Inspectable) error, requestor func(*StreamInputBase)) *SimpleCallbackStreamInput {
	if requestor == nil {
		return nil
	}
	return &SimpleCallbackStreamInput{processor: processor, requestor: requestor}
}

func (s *SimpleCallbackStreamInput) Process(data inspect.Inspectable) error {
	if s.processor == nil {
		return nil
	}
	return s.processor(data)
}

func (s *SimpleCallbackStreamInput) RequestNext() {
	if s.requestor == nil {
		s.RequestData(nil, 0)
	}
	s.requestor(&s.StreamInputBase)
}
