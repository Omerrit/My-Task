package auth

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/registry"
	"gerrit-share.lan/go/actors/replies"
)

type ServiceImpl interface {
	actors.ActorCompatible
	NewConnection(id Id)
	ConnectionClosed(id Id)
	//use Sender() to determine source actor
	RequestUser(connId Id) (Id, error)
	HavePermission(connId Id, cmdName string, path string) (bool, error)
}

type Service struct {
	impl               ServiceImpl
	serviceConnections serviceConnections
	connections        IdSet
}

func (s *Service) onNewConnection(id Id) {
	s.connections.Add(id)
	s.impl.NewConnection(id)
}

func (s *Service) onConnectionClosed(id Id) {
	s.connections.Remove(id)
	base := s.impl.GetBase()
	for actor := range s.serviceConnections.actors[id] {
		base.SendMessage(actor, &connectionClosed{id})
	}
	s.serviceConnections.removeConnection(id)
	s.impl.ConnectionClosed(id)
}

func (s *Service) requestUser(connId Id) (actors.Response, error) {
	if !s.connections.Contains(connId) {
		return nil, ErrNoConnectionId
	}
	id, err := s.impl.RequestUser(connId)
	if err != nil {
		return nil, err
	}
	if id == zeroId {
		return nil, ErrNoUser
	}
	base := s.impl.GetBase()
	sender := base.Sender()
	if !s.serviceConnections.haveService(sender) {
		base.Monitor(sender, func(error) {
			s.serviceConnections.removeService(sender)
		})
	}
	s.serviceConnections.add(sender, connId)
	return &User{id}, nil
}

func (s *Service) getPermission(cmd *getPermission) (replies.BoolResponse, error) {
	if !s.connections.Contains(cmd.connId) {
		return nil, ErrNoConnectionId
	}
	ok, err := s.impl.HavePermission(cmd.connId, cmd.command, cmd.path)
	if err != nil {
		return nil, err
	}
	return replies.Bool(ok), nil
}

func (s *Service) AddBehaviour(actor ServiceImpl, behaviour actors.Behaviour) actors.Behaviour {
	s.impl = actor
	registry.RegisterMe(actor, authServiceName, actor.GetBase().Quit)
	behaviour.AddMessage(new(connectionEstablished), func(message interface{}) {
		s.onNewConnection(message.(*connectionEstablished).Id)
	}).AddMessage(new(connectionClosed), func(message interface{}) {
		s.onConnectionClosed(message.(*connectionClosed).Id)
	})
	behaviour.AddCommand(new(userRequest), func(cmd interface{}) (actors.Response, error) {
		return s.requestUser(cmd.(*userRequest).Id)
	}).AddCommand(new(getPermission), func(cmd interface{}) (actors.Response, error) {
		return s.getPermission(cmd.(*getPermission))
	})
	return behaviour
}

//call this when connection privileges have changed, this will notify everyone interested in it
func (s *Service) PrivilegesChanged(connId Id) {
	if !s.connections.Contains(connId) {
		return
	}
	base := s.impl.GetBase()
	for actor := range s.serviceConnections.actors[connId] {
		base.SendMessage(actor, &connectionPrivilegesChanged{connId})
	}
}
