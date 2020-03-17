package actors

import (
	"gerrit-share.lan/go/common"
)

type ActorSet map[ActorService]common.None

func (a *ActorSet) Add(actor ActorService) {
	if *a == nil {
		*a = make(ActorSet, 1)
	}
	(*a)[actor] = common.None{}
}

func (a ActorSet) Contains(actor ActorService) bool {
	_, ok := a[actor]
	return ok
}

func (a *ActorSet) Remove(actor ActorService) {
	delete(*a, actor)
}

func (a *ActorSet) IsEmpty() bool {
	return len(*a) == 0
}
