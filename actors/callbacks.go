package actors

import (
	"gerrit-share.lan/go/errors"
)

type ActorCallback func(*Actor)

func (a ActorCallback) Call(actor *Actor) {
	if a != nil {
		a(actor)
	}
}

type CommandProcessor func(interface{}) (Response, error)

type MessageProcessor func(interface{})

type PanicProcessor func(errors.StackTraceError)

type StateChangeProcessor func(ActorService, ActorState)
