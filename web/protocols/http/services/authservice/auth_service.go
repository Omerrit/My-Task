package authservice

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/starter"
	"gerrit-share.lan/go/web/auth"
	"log"
)

// just a dummy one atm
type serviceImpl struct {
	actors.Actor
	name    string
	service auth.Service
}

func (s *serviceImpl) NewConnection(id auth.Id) {
}

func (s *serviceImpl) ConnectionClosed(id auth.Id) {
}

func (s *serviceImpl) RequestUser(connId auth.Id) (auth.Id, error) {
	return auth.Id{}, nil
}

func (s *serviceImpl) HavePermission(connId auth.Id, cmdName string, path string) (bool, error) {
	return true, nil
}

func (s *serviceImpl) MakeBehaviour() actors.Behaviour {
	var b actors.Behaviour
	log.Println(s.name, "started")
	var handle starter.Handle
	handle.Acquire(s, handle.DependOn, s.Quit)
	return s.service.AddBehaviour(s, b)
}

func (s *serviceImpl) Shutdown() error {
	log.Println(s.name, "shut down")
	return nil
}

func init() {
	auth.SetAuthCreator(func(s *actors.Actor, name string) (actors.ActorService, error) {
		service := &serviceImpl{name: serviceImplName}
		return s.System().Spawn(service), nil
	})
}
