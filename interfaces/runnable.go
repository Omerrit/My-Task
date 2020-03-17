package interfaces

import (
	"gerrit-share.lan/go/common"
)

type RunnableWithShutdown interface {
	Run(shutdownChannel common.OutSignalChannel) error
	Shutdownable
}

type DummyRunnableWithShutdown common.None

func (DummyRunnableWithShutdown) Run(common.OutSignalChannel) {
}

func (DummyRunnableWithShutdown) Shutdown() error {
	return nil
}

type RunnableWithShutdownFunc func(common.OutSignalChannel) error

func (r RunnableWithShutdownFunc) Run(shutdownChannel common.OutSignalChannel) error {
	return r(shutdownChannel)
}

func (r RunnableWithShutdownFunc) Shutdown() error {
	return nil
}
