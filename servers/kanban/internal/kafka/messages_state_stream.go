package kafka

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
)

type MessagesStream struct {
	messages      Messages
	startOffset   int
	deleteHistory bool
}

func (m *MessagesStream) Add(msg *Message) {
	m.messages = append(m.messages, msg)
}

func (m *MessagesStream) fillArray(array *Messages, offset int, maxLen int) (inspect.Inspectable, int, error) {
	realLen := len(m.messages) - offset + m.startOffset
	if realLen > maxLen {
		realLen = maxLen
	}
	if realLen == 0 {
		return nil, offset, nil
	}
	array.SetLength(realLen)
	nextOffset := offset + copy(*array, m.messages[(offset-m.startOffset):])
	return array, nextOffset, nil
}

func (m *MessagesStream) FillData(data inspect.Inspectable, offset int, maxLen int) (result inspect.Inspectable, nextOffset int, err error) {
	if offset < m.startOffset {
		return data, offset, actors.ErrOffsetOutOfRange
	}
	if maxLen == 0 {
		maxLen = actors.DefaultMaxLen
	}
	if array, ok := data.(*Messages); ok {
		return m.fillArray(array, offset, maxLen)
	} else if _, ok := data.(*Message); ok {
		if (offset - m.startOffset) > len(m.messages) {
			return nil, offset, actors.ErrOffsetOutOfRange
		}
		if (offset - m.startOffset) == len(m.messages) {
			return nil, offset, nil
		}
		return m.messages[offset-m.startOffset], offset + 1, nil
	}
	return m.fillArray(new(Messages), offset, maxLen)
}

func (m *MessagesStream) GetLatestState() (int, actors.DataSource) {
	array := make(Messages, 0, len(m.messages))
	array = append(array, m.messages...)
	return m.startOffset + len(m.messages), &consumerOutput{Messages: array}
}

func (m *MessagesStream) LastOffsetChanged(offset int) {
	if m.deleteHistory && len(m.messages)/2 < (offset-m.startOffset) {
		var messages Messages
		messages.SetLength(len(m.messages) + m.startOffset - offset)
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
