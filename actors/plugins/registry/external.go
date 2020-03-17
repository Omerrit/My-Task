package registry

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/replyprocessors/actorrp"
)

const pluginName = "registry"

func registry(actor actors.ActorCompatible) actors.ActorService {
	return actor.GetBase().System().GetPluginActor(pluginName)
}

func RegisterMe(me actors.ActorCompatible, name string, onError common.ErrorCallback) {
	me.GetBase().SendRequest(registry(me), &registerMe{actorName{name}}, actors.NewErrorProcessor(onError))
}

func RegisterOther(me actors.ActorCompatible, name string, actor actors.ActorService, onError common.ErrorCallback) {
	me.GetBase().SendRequest(actor.System().GetPluginActor(pluginName), &registerOther{actorName{name}, actor}, actors.NewErrorProcessor(onError))
}

func GetActor(me actors.ActorCompatible, name string, onReply actors.ActorServiceCallback, onError common.ErrorCallback) {
	me.GetBase().SendRequest(registry(me), &getActor{actorName{name}}, actorrp.New(onReply, onError))
}

func WaitActor(me actors.ActorCompatible, name string, onReply actors.ActorServiceCallback, onError common.ErrorCallback) {
	me.GetBase().SendRequest(registry(me), &waitActor{actorName{name}}, actorrp.New(onReply, onError))
}

//will not wait for answer if onError is nil (sends message instead of command)
func MarkAsUnregisterable(me actors.ActorCompatible, name string, onError common.ErrorCallback) {
	if onError != nil {
		me.GetBase().SendRequest(registry(me), &willNotRegister{actorName{name}}, actors.NewErrorProcessor(onError))
	} else {
		me.GetBase().SendMessage(registry(me), &willNotRegister{actorName{name}})
	}
}

func Subscribe(me actors.ActorCompatible, input actors.StreamInput, onError common.ErrorCallback) {
	me.GetBase().RequestStream(input, registry(me), new(subscribe), onError)
}

func init() {
	actors.AddPlugin(pluginName, func() actors.BehavioralActor { return new(registryActor) })

}
