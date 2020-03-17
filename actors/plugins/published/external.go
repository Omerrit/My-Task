package published

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/common"
)

const pluginName = "published"

func published(actor actors.ActorCompatible) actors.ActorService {
	return actor.GetBase().System().GetPluginActor(pluginName)
}

func Publish(me actors.ActorCompatible, onError common.ErrorCallback) {
	me.GetBase().SendRequest(published(me), &publish{me.GetBase().Service()}, actors.NewErrorProcessor(onError))
}

func Subscribe(me actors.ActorCompatible, input actors.StreamInput, onError common.ErrorCallback) {
	me.GetBase().RequestStream(input, published(me), new(subscribe), onError)
}
