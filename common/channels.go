package common

import "context"

type (
	None             struct{}
	SignalChannel    chan struct{}
	InSignalChannel  chan<- struct{}
	OutSignalChannel <-chan struct{}
)

func (s SignalChannel) ToInput() InSignalChannel {
	return InSignalChannel((chan struct{})(s))
}

func (s SignalChannel) ToOutput() OutSignalChannel {
	return OutSignalChannel((chan struct{})(s))
}

func (i OutSignalChannel) Await() bool {
	_, ok := <-i
	return ok
}

func (i OutSignalChannel) AwaitWithShutdown(shutdown OutSignalChannel) bool {
	select {
	case <-i:
		return true
	case <-shutdown:
		return false
	}
}

func (i OutSignalChannel) AwaitWithContext(ctx context.Context) bool {
	if ctx == nil {
		return i.Await()
	}
	return i.AwaitWithShutdown(ctx.Done())
}

func WhenDone(f func()) OutSignalChannel {
	signalChannel := make(SignalChannel)
	go func() {
		defer close(signalChannel)
		f()
	}()
	return signalChannel.ToOutput()
}
