package actors

import (
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/inspect"
	"math/big"
)

type actorResponseVisitor Actor

func (a *actorResponseVisitor) Reply(data inspect.Inspectable) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyInt(data int) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyInt32(data int32) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyInt64(data int64) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyFloat32(data float32) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyFloat64(data float64) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyString(data string) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyBytes(data []byte) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyBool(data bool) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyBigInt(data *big.Int) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyRat(data *big.Rat) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) ReplyBigFloat(data *big.Float) {
	(*Actor)(a).reply(data)
}

func (a *actorResponseVisitor) MakeResponsePromise(onCancel common.SimpleCallback) ResponsePromise {
	if a.currentCommand.isValid() {
		promise := ResponsePromise{a.currentCommand.promiseId, (*Actor)(a), nil}
		if !a.activePromises.Contains(promise.id) { //we're not reprocessing command
			a.currentCommand.preReply((*Actor)(a).Service())
		}
		a.activePromises.Add(promise.id, onCancel)
		a.currentCommand.invalidate()
		return promise
	}
	return ResponsePromise{}
}
