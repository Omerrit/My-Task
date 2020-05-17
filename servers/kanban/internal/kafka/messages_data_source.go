package kafka

import (
	"gerrit-share.lan/go/actors"
	"log"
)

type MessageDataSource struct {
	consumerOutput
	maxOffset          int64
	lastOffset         int64
	onNewDataAvailable func()
}

func (m *MessageDataSource) Init(lastOffset int64) func(*Message) (error, bool) {
	m.lastOffset = lastOffset
	return m.push
}

func (m *MessageDataSource) push(message *Message) (error, bool) {
	log.Printf("[%p] got offsets %d need max offset %d\n", m, message.Offset, m.maxOffset)
	if message.Offset > m.maxOffset {
		return nil, true
	}
	if message.Offset < 0 {
		return nil, m.lastOffset >= m.maxOffset
	}
	m.Messages = append(m.Messages, message)
	m.lastOffset = message.Offset
	m.onNewDataAvailable()
	if message.Offset == m.maxOffset {
		m.Messages = append(m.Messages, &Message{nil, nil, m.maxOffset})
		return nil, true
	}
	return nil, false
}

func (m *MessageDataSource) RequestNext(onDataAvailable func(), maxOffset int64) bool {
	if maxOffset <= m.lastOffset {
		log.Printf("[%p] requesting offset %d (no data to give)\n", m, maxOffset)
		return false
	}
	log.Printf("[%p] requesting offset %d\n", m, maxOffset)
	m.onNewDataAvailable = onDataAvailable
	m.maxOffset = maxOffset
	return true
}

func init() {
	var _ actors.DynamicDataSource = (*MessageDataSource)(nil)
}
