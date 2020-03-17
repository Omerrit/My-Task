package interfaces

import (
	"context"
	"gerrit-share.lan/go/common"
)

type (
	BytesPromise interface {
		Awaitable
		Get() []byte
	}

	StringPromise interface {
		Awaitable
		Get() []string
	}

	StringOutChannelPromise interface {
		Awaitable
		Get() <-chan string
	}

	UntypedPromise interface {
		Awaitable
		Get() (interface{}, error)
	}
)

type AwaitableArray []Awaitable

func (a *AwaitableArray) Add(awaitable Awaitable) {
	if awaitable != nil {
		*a = append(*a, awaitable)
	}
}

func (a *AwaitableArray) Compact() Awaitable {
	switch len(*a) {
	case 0:
		return nil
	case 1:
		return (*a)[0]
	default:
		return a
	}
}

func (a *AwaitableArray) AwaitWithShutdown(shutdown common.OutSignalChannel) bool {
	result := true
	for _, awaitable := range *a {
		if !awaitable.AwaitWithShutdown(shutdown) {
			result = false
		}
	}
	return result
}

func (a *AwaitableArray) AwaitWithContext(ctx context.Context) bool {
	result := true
	for _, awaitable := range *a {
		if !awaitable.AwaitWithContext(ctx) {
			result = false
		}
	}
	return result
}

func (a *AwaitableArray) Await() bool {
	result := true
	for _, awaitable := range *a {
		if !awaitable.Await() {
			result = false
		}
	}
	return result
}

func (a *AwaitableArray) Clear() {
	*a = nil
}
