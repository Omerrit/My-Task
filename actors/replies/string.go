package replies

import "gerrit-share.lan/go/actors"

type String string

func (s String) Visit(visitor actors.ResponseVisitor) {
	visitor.ReplyString(string(s))
}

func (String) Deliver(string) {}

type StringPromise struct {
	actors.ResponsePromise
}

func (s *StringPromise) Deliver(value string) {
	s.DeliverString(value)
}

type StringResponse interface {
	actors.Response
	Deliver(string)
}
