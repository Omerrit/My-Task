package actors

import (
	"gerrit-share.lan/go/actors/internal/queue"
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"gerrit-share.lan/go/interfaces"
	"math/big"
)

type ActorService interface {
	interfaces.NamedClosableService
	inspect.Inspectable
	Response
	SendMessage(inspect.Inspectable)
	SendQuit(error)
	CloseError() error
	System() *System

	enqueue(message interface{})
	start(*System)
	getActor() ActorCompatible
	init(ActorCompatible, *queue.Queue)
}

type ActorCompatible interface {
	GetBase() *Actor
	Run() error
}

type ResponseVisitor interface {
	ReplyInt(reply int)
	ReplyInt32(reply int32)
	ReplyInt64(reply int64)
	ReplyFloat32(reply float32)
	ReplyFloat64(reply float64)
	ReplyString(reply string)
	ReplyBytes(reply []byte)
	ReplyBool(reply bool)
	ReplyBigInt(reply *big.Int)
	ReplyRat(reply *big.Rat)
	ReplyBigFloat(reply *big.Float)
	Reply(reply inspect.Inspectable)
	MakeResponsePromise(onCancel common.SimpleCallback) ResponsePromise
}

type TypeSystem interface {
	Creator(string) inspectables.Creator
}

type ActorServiceCallback func(ActorService)

func (a ActorServiceCallback) Call(service ActorService) {
	if a == nil {
		return
	}
	a(service)
}
