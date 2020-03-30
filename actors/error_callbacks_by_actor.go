package actors

import (
	"gerrit-share.lan/go/common"
)

type ActorErrorCallbacks map[ActorService][]common.ErrorCallback

func (a *ActorErrorCallbacks) Add(service ActorService, callback common.ErrorCallback) {
	if *a == nil {
		*a = make(ActorErrorCallbacks, 1)
	}
	(*a)[service] = append((*a)[service], callback)
}

func (a *ActorErrorCallbacks) Call(service ActorService, err error) {
	for _, callback := range (*a)[service] {
		callback.Call(err)
	}
}

func (a *ActorErrorCallbacks) CallAndRemove(service ActorService, err error) {
	callbacks, ok := (*a)[service]
	if !ok {
		return
	}
	for _, callback := range callbacks {
		callback.Call(err)
	}
	delete(*a, service)
}

func (a *ActorErrorCallbacks) IsEmpty() bool {
	return len(*a) == 0
}
