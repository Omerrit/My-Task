package actors

import (
	"sync"
)

//this should initialize actor implementation, not spawn the actor
//system is not available inside this call, use MakeBehaviour for actual initialization
type ActorCreator func() BehavioralActor

var plugins sync.Map

//Adds plugin, replaces on name collision.
//If creator is nil then remove plugin if present and mark plugin name as unavailable.
//It's not possible to add plugin with a name that was marked as anavailable
func AddPlugin(name string, creator ActorCreator) {
	value, ok := plugins.Load(name)
	if ok && value.(ActorCreator) == nil {
		return
	}
	plugins.Store(name, creator)
}

func iteratePlugins(iterator func(string, ActorCreator) bool) {
	plugins.Range(func(key interface{}, value interface{}) bool {
		return iterator(key.(string), value.(ActorCreator))
	})
}
