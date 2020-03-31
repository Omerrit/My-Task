package auth

import (
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/registry"
	"gerrit-share.lan/go/actors/starter"
)

func SetAuthCreator(creator starter.ServiceCreator) {
	starter.SetCreator(authServiceName, creator)
}

func init() {
	starter.SetCreatorIfNotPresent(authServiceName,
		func(parent *actors.Actor, name string) (actors.ActorService, error) {
			registry.MarkAsUnregisterable(parent, authServiceName, nil)
			return nil, fmt.Errorf("%w: auth service (%s) will not start", actors.ErrNotGonnaHappen, authServiceName)
		})
}
