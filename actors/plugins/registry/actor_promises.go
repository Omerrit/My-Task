package registry

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/common"
)

type actorPromises struct {
	promises map[string]map[*actors.ResponsePromise]common.None
	names    map[*actors.ResponsePromise]string
}

func (a *actorPromises) add(name string) *actors.ResponsePromise {
	if a.promises == nil {
		a.promises = make(map[string]map[*actors.ResponsePromise]common.None, 1)
	}
	promises := a.promises[name]
	if promises == nil {
		promises = make(map[*actors.ResponsePromise]common.None, 1)
	}
	result := new(actors.ResponsePromise)
	promises[result] = common.None{}
	a.promises[name] = promises
	if a.names == nil {
		a.names = make(map[*actors.ResponsePromise]string, 1)
	}
	a.names[result] = name
	return result
}

func (a *actorPromises) remove(promise *actors.ResponsePromise) {
	name, ok := a.names[promise]
	if !ok {
		return
	}
	promises := a.promises[name]
	delete(promises, promise)
	a.promises[name] = promises
	delete(a.names, promise)
}

func (a *actorPromises) fulfill(name string, actor actors.ActorService) {
	for promise := range a.promises[name] {
		promise.DeliverUntyped(actor)
		delete(a.names, promise)
	}
	delete(a.promises, name)
}

func (a *actorPromises) fail(name string) {
	for promise := range a.promises[name] {
		promise.Fail(actors.ErrNotGonnaHappen)
		delete(a.names, promise)
	}
	delete(a.promises, name)
}
