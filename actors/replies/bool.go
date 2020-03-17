package replies

import (
	"gerrit-share.lan/go/actors"
)

type Bool bool

func (b Bool) Visit(visitor actors.ResponseVisitor) {
	visitor.ReplyBool(bool(b))
}

func (Bool) Deliver(bool) {}

type BoolPromise struct {
	actors.ResponsePromise
}

func (b *BoolPromise) Deliver(value bool) {
	b.DeliverBool(value)
}

type BoolResponse interface {
	actors.Response
	Deliver(bool)
}
