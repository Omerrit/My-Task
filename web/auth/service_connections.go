package auth

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/common"
)

type IdSet map[Id]common.None

func (s *IdSet) Add(id Id) {
	if *s == nil {
		*s = make(IdSet, 1)
	}
	(*s)[id] = common.None{}
}

func (s *IdSet) Remove(id Id) {
	delete(*s, id)
}

func (s *IdSet) Contains(id Id) bool {
	_, ok := (*s)[id]
	return ok
}

func (s *IdSet) IsEmpty() bool {
	return len(*s) == 0
}

type serviceConnections struct {
	connections map[actors.ActorService]IdSet
	actors      map[Id]actors.ActorSet
}

func (s *serviceConnections) add(actor actors.ActorService, connId Id) {
	if s.actors == nil {
		s.actors = make(map[Id]actors.ActorSet, 1)
	}
	actorSet := s.actors[connId]
	actorSet.Add(actor)
	s.actors[connId] = actorSet
	if s.connections == nil {
		s.connections = make(map[actors.ActorService]IdSet, 1)
	}
	connections := s.connections[actor]
	connections.Add(connId)
	s.connections[actor] = connections
}

func (s *serviceConnections) haveService(actor actors.ActorService) bool {
	_, ok := s.connections[actor]
	return ok
}

func (s *serviceConnections) removeConnection(connId Id) {
	for actor := range s.actors[connId] {
		connections := s.connections[actor]
		connections.Remove(connId)
		if connections.IsEmpty() {
			delete(s.connections, actor)
		} else {
			s.connections[actor] = connections
		}
	}
	delete(s.actors, connId)
}

func (s *serviceConnections) removeService(service actors.ActorService) {
	for connId := range s.connections[service] {
		actors := s.actors[connId]
		actors.Remove(service)
		if actors.IsEmpty() {
			delete(s.actors, connId)
		} else {
			s.actors[connId] = actors
		}
	}
	delete(s.connections, service)
}
