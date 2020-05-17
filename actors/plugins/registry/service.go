package registry

import (
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/utils/sets"
	"log"
)

type registryActor struct {
	actors.Actor
	registry       map[actors.ActorService]string
	names          infoStateChangeStream
	broadcaster    actors.StateBroadcaster
	promises       actorPromises
	forbiddenNames sets.String
}

func (r *registryActor) Shutdown() error {
	log.Println("registry plugin shut down")
	return nil
}

func (r *registryActor) register(name string, service actors.ActorService) error {
	if r.forbiddenNames.Contains(name) {
		return actors.ErrNotGonnaHappen
	}
	if service, ok := r.names.Get(name); ok {
		delete(r.registry, service)
	}
	r.names.Add(name, service)
	r.registry[service] = name
	r.promises.fulfill(name, service)
	r.broadcaster.NewDataAvailable()
	r.Monitor(service, func(error) {
		name, ok := r.registry[service]
		if !ok {
			return
		}
		delete(r.registry, service)
		r.names.Remove(name)
	})
	return nil
}

func (r *registryActor) getActor(name string) (actors.ActorService, error) {
	if r.forbiddenNames.Contains(name) {
		return nil, fmt.Errorf("%w: %s will never start", actors.ErrNotGonnaHappen, name)
	}
	actor, ok := r.names.Get(name)
	if !ok {
		return nil, actors.ErrNotFound
	}
	return actor, nil
}

func (r *registryActor) waitActor(name string) (actors.Response, error) {
	if r.forbiddenNames.Contains(name) {
		return nil, fmt.Errorf("%w: %s will never start", actors.ErrNotGonnaHappen, name)
	}
	actor, ok := r.names.Get(name)
	if !ok {
		promise := r.promises.add(name)
		promise.OnCancel(func() {
			r.promises.remove(promise)
		})
		return promise, nil
	}
	return actor, nil
}

func (r *registryActor) willNotRegister(name string) error {
	if r.names.Contains(name) {
		return actors.ErrAlreadyRegistered
	}
	r.forbiddenNames.Add(name)
	return nil
}

func (r *registryActor) subscribe(cmd *subscribe) {
	r.InitStreamOutput(r.broadcaster.AddOutput(), cmd)
}

func (r *registryActor) MakeBehaviour() actors.Behaviour {
	log.Println("registry plugin started")
	r.registry = make(map[actors.ActorService]string)
	r.names.names = make(map[string]actors.ActorService)
	r.broadcaster = actors.NewBroadcaster(&r.names)
	r.broadcaster.CloseWhenActorCloses()
	var b actors.Behaviour
	b.AddCommand(new(registerMe), func(cmd interface{}) (actors.Response, error) {
		return nil, r.register(cmd.(*registerMe).name, r.Sender())
	}).AddCommand(new(registerOther), func(cmd interface{}) (actors.Response, error) {
		command := cmd.(*registerOther)
		return nil, r.register(command.name, command.actor)
	}).AddCommand(new(getActor), func(cmd interface{}) (actors.Response, error) {
		return r.getActor(cmd.(*getActor).name)
	}).AddCommand(new(waitActor), func(cmd interface{}) (actors.Response, error) {
		return r.waitActor(cmd.(*waitActor).name)
	}).AddCommand(new(willNotRegister), func(cmd interface{}) (actors.Response, error) {
		return nil, r.willNotRegister(cmd.(*willNotRegister).name)
	}).AddCommand(new(subscribe), func(cmd interface{}) (actors.Response, error) {
		r.subscribe(cmd.(*subscribe))
		return nil, nil
	}).Result(new(InfoArray))
	b.AddMessage(new(willNotRegister), func(msg interface{}) {
		r.willNotRegister(msg.(*willNotRegister).name)
	})
	b.Name = "actor registry plugin"
	return b
}
