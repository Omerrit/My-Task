package replies

import (
	"gerrit-share.lan/go/actors"
)

type BytesPromise struct {
	actors.ResponsePromise
}

func (b *BytesPromise) Deliver(value []byte) {
	b.DeliverBytes(value)
}
