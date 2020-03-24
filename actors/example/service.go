package example

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/registry"
	"gerrit-share.lan/go/actors/starter"
)

type service struct {
	actors.Actor
}

func (s *service) processMessage(msg *message) {
	//process message here
}

//command processor, don't use untyped inputs or results here
func (s *service) getStuff(msg *getStuff) *stuff {
	return &stuff{name: msg.path, value: 42}
}

//mandatory initialization function, runs in actor's own goroutine right after it starts
func (s *service) MakeBehaviour() actors.Behaviour {
	//Do initialization that requires a running actor here
	registry.RegisterMe(s, serviceName, s.Quit)
	//after that create behaviour
	var behaviour actors.Behaviour
	//use lambdas as message and command processors that call actual processor with typed arguments and results
	behaviour.AddMessage(new(message), func(msg interface{}) {
		//it's guaranteed that the handler receives only messages of the real type it expects
		s.processMessage(msg.(*message))
	})
	behaviour.AddCommand(new(getStuff), func(cmd interface{}) (actors.Response, error) {
		return s.getStuff(cmd.(*getStuff)), nil
	})
	return behaviour
}

func init() {
	//make the service launch automatically in starter
	starter.SetCreator(serviceName, func(parent *actors.Actor, name string) (actors.ActorService, error) {
		srv := new(service)
		//parent is a service that will be killed on ctrl-c, kill our actor when that happens
		srv.DependOn(parent.Service())
		return parent.System().Spawn(srv), nil
	})
}
