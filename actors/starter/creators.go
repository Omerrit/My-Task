package starter

import (
	"gerrit-share.lan/go/actors/services/autorestarter"
)

type ServiceCreator autorestarter.ServiceMaker

var defaultServiceCreators = make(map[string]ServiceCreator)

func SetCreator(name string, creator ServiceCreator) {
	if creator == nil {
		delete(defaultServiceCreators, name)
	}
	defaultServiceCreators[name] = creator
}

func SetCreatorIfNotPresent(name string, creator ServiceCreator) {
	_, ok := defaultServiceCreators[name]
	if ok {
		return
	}
	SetCreator(name, creator)
}
