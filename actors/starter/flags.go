package starter

import (
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/utils/maps"
)

var serviceFlagInitializers maps.SimpleCallbacksForNames
var serviceFlagProcessors maps.SimpleCallbacksForNames

func SetFlagInitializer(name string, initializer common.SimpleCallback) {
	if initializer == nil {
		serviceFlagInitializers.Delete(name)
	} else {
		serviceFlagInitializers.Add(name, initializer)
	}
}

func SetFlagProcessor(name string, processor common.SimpleCallback) {
	if processor == nil {
		serviceFlagProcessors.Delete(name)
	} else {
		serviceFlagProcessors.Add(name, processor)
	}
}
