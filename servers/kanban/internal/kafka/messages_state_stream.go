package kafka

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
)

type MessagesStream struct {
	messages        Messages
	startOffset     int64
	haveSubscribers bool
}

func (m *MessagesStream) StartOffset() int64 {
	return m.startOffset
}

func (m *MessagesStream) Add(msg *Message) {
	if msg.Offset < 0 {
		return
	}
	if m.haveSubscribers {
		m.messages = append(m.messages, msg)
	} else {
		if msg.Offset >= 0 {
			m.startOffset = msg.Offset
		}
	}
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
		return m.messages[int(offset-m.startOffset)], offset + 1, nil
	}
	return m.fillArray(new(Messages), offset, maxLen)
}

func (m *MessagesStream) GetLatestState() (int64, actors.DataSource) {
	m.haveSubscribers = true
	return m.startOffset + int64(len(m.messages)), &MessageDataSource{}
}

func (m *MessagesStream) LastOffsetChanged(offset int64) {
	if len(m.messages)/2 < int(offset-m.startOffset) {
		var messages Messages
		messages.SetLength(len(m.messages) - int(offset-m.startOffset))
		copy(messages, m.messages[int(offset-m.startOffset):])
		m.messages = messages
		m.startOffset = offset
	}
}

func (m *MessagesStream) NoMoreSubscribers() {
	m.haveSubscribers = false
	if len(m.messages) == 0 {
		return
	}
	m.startOffset = m.messages[len(m.messages)-1].Offset
	m.messages = m.messages[:0]
}
