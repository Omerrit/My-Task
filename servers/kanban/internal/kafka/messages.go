package kafka

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"github.com/Shopify/sarama"
)

type Message struct {
	Key    []byte
	Value  []byte
	Offset int64
}

const OffsetUninitialized = -10

func newMessage(msg *sarama.ConsumerMessage) *Message {
	if msg == nil {
		return nil
	}
	return &Message{msg.Key, msg.Value, msg.Offset}
}

const messageName = packageName + "message"

func (m *Message) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(messageName, "message")
	{
		objectInspector.ByteString(&m.Key, "key", true, "key of message")
		objectInspector.ByteString(&m.Value, "value", true, "value of message")
		objectInspector.Int64(&m.Offset, "offset", true, "message offset")
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(messageName, func() inspect.Inspectable { return new(Message) })
}

type Messages []*Message

func (m *Messages) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(packageName+".messages", messageName, "array of messages")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*m))
		} else {
			m.SetLength(arrayInspector.GetLength())
		}
		for index := range *m {
			(*m)[index].Inspect(arrayInspector.Value())
		}
		arrayInspector.End()
	}
}

func (m *Messages) SetLength(length int) {
	if cap(*m) < length {
		*m = make(Messages, length)
	} else {
		*m = (*m)[:length]
	}
}

func (m *Messages) Add(msg *sarama.ConsumerMessage) {
	*m = append(*m, newMessage(msg))
}
