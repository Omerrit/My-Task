package kanban

import (
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/servers/kanban/internal/kafka"
	"net/http"
)

type streamInput struct {
	actors.StreamInputBase
	writer  http.ResponseWriter
	flusher http.Flusher
}

func newStreamInput(writer http.ResponseWriter) *streamInput {
	return &streamInput{writer: writer, flusher: writer.(http.Flusher)}
}

func (s *streamInput) Process(data inspect.Inspectable) error {
	msgs := data.(*kafka.Messages)
	for _, msg := range *msgs {
		err := s.processMessage(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *streamInput) RequestNext() {
	s.StreamInputBase.RequestData(new(kafka.Messages), 10)
}

func (s *streamInput) processMessage(message *kafka.Message) error {
	fmt.Fprintf(s.writer, "id: %d\n", message.Offset)
	fmt.Fprintf(s.writer, "data: %s\n", message.Key)
	fmt.Fprintf(s.writer, "data: %s\n\n", message.Value)
	s.flusher.Flush()
	return nil
}
