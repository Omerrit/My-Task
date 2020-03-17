package actors

import (
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/inspect"
	"math/big"
)

type ResponsePromise struct {
	id       promiseId
	actor    *Actor
	onCancel common.SimpleCallback
}

func (r *ResponsePromise) Visit(visitor ResponseVisitor) {
	*r = visitor.MakeResponsePromise(r.onCancel)
}

func (r *ResponsePromise) OnCancel(callback common.SimpleCallback) {
	r.onCancel = callback
}

func (r *ResponsePromise) Cancel(err error) {
	if r.actor == nil {
		return
	}
	r.actor.cancelCommandProcessingFromPromise(r.id, err)
	r.actor = nil
}

func (r *ResponsePromise) Fail(err error) {
	if r.actor == nil {
		return
	}
	r.actor.replyWithErrorFromPromise(r.id, err)
	r.actor = nil
}

func (r *ResponsePromise) DeliverUntyped(data inspect.Inspectable) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverInt(data int) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverInt32(data int32) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverInt64(data int64) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverFloat32(data float32) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverFloat64(data float64) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverString(data string) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverBytes(data []byte) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverBool(data bool) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverBigInt(data *big.Int) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverRat(data *big.Rat) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}

func (r *ResponsePromise) DeliverBigFloat(data *big.Float) {
	if r.actor == nil {
		return
	}
	r.actor.replyFromPromise(r.id, data)
	r.actor = nil
}
