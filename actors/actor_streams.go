package actors

import ()

type streamInputState int

const (
	streamInputInitialized streamInputState = iota
	streamInputAwaitingData
)

type streamInputs map[streamId]StreamInput

func (s *streamInputs) Add(id streamId, input StreamInput) {
	if *s == nil {
		*s = make(streamInputs, 1)
	}
	(*s)[id] = input
}

func (s *streamInputs) Remove(id streamId) {
	delete(*s, id)
}

func (s *streamInputs) Clear() {
	*s = make(streamInputs, 1)
}

func (s *streamInputs) IsEmpty() bool {
	return len(*s) == 0
}

type streamOutInfo struct {
	output      StreamOutput
	dataRequest streamRequest
}

type streamOutputs map[outputId]streamOutInfo

func (s *streamOutputs) Add(id outputId, output StreamOutput) {
	if *s == nil {
		*s = make(streamOutputs, 1)
	}
	(*s)[id] = streamOutInfo{output, streamRequest{}}
}

func (s *streamOutputs) Remove(id outputId) {
	delete(*s, id)
}

func (s *streamOutputs) Clear() {
	*s = make(streamOutputs, 1)
}

func (s *streamOutputs) IsEmpty() bool {
	return len(*s) == 0
}
