package actors

import (
	"gerrit-share.lan/go/common"
)

type promiseIdCallbacks map[promiseId]common.SimpleCallback

func (p *promiseIdCallbacks) Add(id promiseId, callback common.SimpleCallback) {
	if *p == nil {
		*p = make(promiseIdCallbacks, 1)
	}
	(*p)[id] = callback
}

func (p *promiseIdCallbacks) Delete(id promiseId) {
	delete(*p, id)
}

func (p promiseIdCallbacks) Contains(id promiseId) bool {
	_, ok := p[id]
	return ok
}

func (p *promiseIdCallbacks) Clear() {
	for key := range *p {
		delete(*p, key)
	}
}

func (p promiseIdCallbacks) IsEmpty() bool {
	return len(p) == 0
}

type promiseIdSet map[promiseId]common.None

func (p *promiseIdSet) Add(id promiseId) {
	if *p == nil {
		*p = make(promiseIdSet, 1)
	}
	(*p)[id] = common.None{}
}

func (p *promiseIdSet) Delete(id promiseId) {
	delete(*p, id)
}

func (p promiseIdSet) Contains(id promiseId) bool {
	_, ok := p[id]
	return ok
}

func (p *promiseIdSet) Clear() {
	for key := range *p {
		delete(*p, key)
	}
}

func (p promiseIdSet) IsEmpty() bool {
	return len(p) == 0
}
