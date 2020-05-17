package actors

import ()

type BehaviourMaker func(*Actor) Behaviour

type simpleActor struct {
	Actor
	makeBehaviour func(*Actor) Behaviour
}

func NewSimpleActor(behaviourMaker BehaviourMaker) BehavioralActor {
	return &simpleActor{makeBehaviour: behaviourMaker}
}

func (s *simpleActor) MakeBehaviour() Behaviour {
	if s.makeBehaviour == nil {
		return Behaviour{}
	}
	return s.makeBehaviour(&s.Actor)
}

type SimpleInitializer func(*Actor)

type simpleInitializerActor struct {
	Actor
	initializer SimpleInitializer
}

func (s *simpleInitializerActor) MakeBehaviour() Behaviour {
	if s.initializer == nil {
		return Behaviour{}
	}
	s.initializer(s.GetBase())
	return Behaviour{}
}

func NewSimpleInitializerActor(simpleInitializer SimpleInitializer) BehavioralActor {
	return &simpleInitializerActor{initializer: simpleInitializer}
}

type simpleNamedInitializerActor struct {
	Actor
	initializer SimpleInitializer
	name        string
}

func (s *simpleNamedInitializerActor) MakeBehaviour() Behaviour {
	behaviour := Behaviour{Name: s.name}
	if s.initializer == nil {
		return behaviour
	}
	s.initializer(s.GetBase())
	return behaviour
}

func NewSimpleNamedInitializerActor(name string, simpleInitializer SimpleInitializer) BehavioralActor {
	return &simpleNamedInitializerActor{initializer: simpleInitializer, name: name}
}
