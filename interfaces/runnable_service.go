package interfaces

import ()

type (
	RunnableServicePromise interface {
		Awaitable
		Get() ClosableService
	}

	ServiceRegistry interface {
		ClosableService
		GetService(name string, requester ClosableService) RunnableServicePromise
		PutService(service ClosableServiceStarter)
		RemoveService(name string)
		RemoveServiceByName(name string)
		ForceStartService(name string)
	}

	//not an actual service
	ClosableServiceStarter interface {
		Named
		//may block for registry answers
		Start(registry ServiceRegistry) ClosableService
	}
)
