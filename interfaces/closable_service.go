package interfaces

import (
	"context"
	"gerrit-share.lan/go/common"
)

type (
	ClosableService interface {
		Close() error                         //cleanly shut down
		DoneChannel() common.OutSignalChannel //channel that closes when service shuts down
		Shutdown()
	}

	CloserWithContext interface {
		CloseWithContext(context.Context) error
	}

	NamedClosableService interface {
		ClosableService
		Named
	}

	Initializable interface {
		InitializedChannel() common.OutSignalChannel
	}

	InitializableService interface {
		ClosableService
		Initializable
	}

	ServiceHolder interface {
		Service() ClosableService
	}
)

type DummyClosableService common.None

func (DummyClosableService) Close() error {
	return nil
}

func (DummyClosableService) DoneChannel() common.OutSignalChannel {
	result := make(common.SignalChannel)
	close(result)
	return result.ToOutput()
}

type DummyNamedClosableService struct {
	DummyClosableService
	DummyNamed
}

type DummyInitializable common.None

func (DummyInitializable) InitializedChannel() common.OutSignalChannel {
	result := make(common.SignalChannel)
	close(result)
	return result.ToOutput()
}

type DummyInitializableService struct {
	DummyClosableService
	DummyInitializable
}
