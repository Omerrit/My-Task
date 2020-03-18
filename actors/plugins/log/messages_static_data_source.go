package log

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
)

type staticMessagesDataSource struct {
	Messages      Messages
	currentOffset int
}

func (s *staticMessagesDataSource) Acknowledged() {
	s.Messages = s.Messages[s.currentOffset:]
	s.currentOffset = 0
}

func (s *staticMessagesDataSource) fillArray(array *Messages, maxLen int) inspect.Inspectable {
	realLen := len(s.Messages)
	if realLen > maxLen {
		realLen = maxLen
	}
	if realLen == 0 {
		return nil
	}
	array.SetLength(realLen)
	copy(*array, s.Messages)
	if realLen > s.currentOffset {
		s.currentOffset = realLen
	}
	return array
}

func (s *staticMessagesDataSource) FillData(data inspect.Inspectable, maxLen int) (inspect.Inspectable, error) {
	if len(s.Messages) == 0 {
		return nil, nil
	}
	if maxLen == 0 {
		maxLen = actors.DefaultMaxLen
	}
	if array, ok := data.(*Messages); ok {
		return s.fillArray(array, maxLen), nil
	} else if value, ok := data.(*Message); ok {
		if len(s.Messages) == 0 {
			return nil, nil
		}
		*value = s.Messages[0]
		if s.currentOffset == 0 {
			s.currentOffset = 1
		}
		return value, nil
	} else if data == nil {
		return s.fillArray(new(Messages), maxLen), nil
	}
	return nil, actors.ErrWrongTypeRequested
}

func init() {
	var _ actors.DataSource = (*staticMessagesDataSource)(nil)
}
