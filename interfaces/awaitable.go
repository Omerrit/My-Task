package interfaces

import (
	"context"
	"gerrit-share.lan/go/common"
)

type Awaitable interface {
	Await() bool
	AwaitWithShutdown(common.OutSignalChannel) bool
	AwaitWithContext(ctx context.Context) bool
}

type DummyAwaitable common.None

func (DummyAwaitable) Await() bool {
	return true
}

func (DummyAwaitable) AwaitWithShutdown(common.OutSignalChannel) bool {
	return true
}

type AwaitableFunc func(common.OutSignalChannel) bool

func (f AwaitableFunc) Await() bool {
	return f(nil)
}

func (f AwaitableFunc) AwaitWithShutdown(shutdown common.OutSignalChannel) bool {
	return f(shutdown)
}

func (f AwaitableFunc) AwaitWithContext(ctx context.Context) bool {
	if ctx == nil {
		return f(nil)
	}
	return f(common.OutSignalChannel(ctx.Done()))
}
