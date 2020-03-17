package interfaces

import ()

type ClosableServiceCallback func(ClosableService)

func (c ClosableServiceCallback) Call(service ClosableService) {
	if c != nil {
		c(service)
	}
}
