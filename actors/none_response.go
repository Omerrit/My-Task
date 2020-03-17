package actors

import (
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/inspect"
)

type NoneResponse interface {
	Response
	Deliver()
}

type NoneReply common.None

func (n NoneReply) Visit(visitor ResponseVisitor) {
	visitor.Reply(inspect.DummyInspectable(common.None{}))
}

func (n NoneReply) Deliver() {

}

type NonePromise struct {
	ResponsePromise
}

func (n *NonePromise) Deliver() {
	n.DeliverUntyped(inspect.DummyInspectable(common.None{}))
}
