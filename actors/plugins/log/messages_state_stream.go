package log

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
)

type MessagesStream struct {
	messages      Messages
	startOffset   int64
	deleteHistory bool
}

func (m *MessagesStream) Add(msg Message) {
	m.messages = append(m.messages, msg)
}

func (m *MessagesStream) fillArray(array *Messages, offset int64, maxLen int) (inspect.Inspectable, int64, error) {
	realLen := len(m.messages) - int(offset-m.startOffset)
	if realLen > maxLen {
		realLen = maxLen
	}
	if realLen == 0 {
		return nil, offset, nil
	}
	array.SetLength(realLen)
	nextOffset := offset + int64(copy(*array, m.messages[int(offset-m.startOffset):]))
	return array, nextOffset, nil
}

func (m *MessagesStream) FillData(data inspect.Inspectable, offset int64, maxLen int) (result inspect.Inspectable, nextOffset int64, err error) {
	if offset < m.startOffset {
		return data, offset, actors.ErrOffsetOutOfRange
	}
	if maxLen == 0 {
		maxLen = actors.DefaultMaxLen
	}
	if array, ok := data.(*Messages); ok {
		return m.fillArray(array, offset, maxLen)
	} else if _, ok := data.(*Message); ok {
		if int(offset-m.startOffset) > len(m.messages) {
			return nil, offset, actors.ErrOffsetOutOfRange
		}
		if int(offset-m.startOffset) == len(m.messages) {
			return nil, offset, nil
		}
		return &m.messages[int(offset-m.startOffset)], offset + 1, nil
	}
	return m.fillArray(new(Messages), offset, maxLen)
}

func (m *MessagesStream) GetLatestState() (int64, actors.DataSource) {
	array := make(Messages, 0, len(m.messages))
	array = append(array, m.messages...)
	return m.startOffset + int64(len(m.messages)), &staticMessagesDataSource{Messages: array}
}

func (m *MessagesStream) LastOffsetChanged(offset int64) {
	if m.deleteHistory && len(m.messages)/2 < int(offset-m.startOffset) {
		var messages Messages
		messages.SetLength(len(m.messages) - int(offset-m.startOffset))
		copy(messages, m.messages[offset-m.startOffset:])
		m.messages = messages
		m.startOffset = offset
	}
}

func (m *MessagesStream) NoMoreSubscribers() {
	if m.deleteHistory {
		m.messages = m.messages[:0]
		m.startOffset = 0
	}
}
