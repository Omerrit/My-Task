package actors

import (
	"gerrit-share.lan/go/common"
)

type Runner func(*Actor) error

type runnerActor struct {
	Actor
	runner Runner
}

func (s *runnerActor) Run() error {
	if s.runner == nil {
		return nil
	}
	return s.runner(s.GetBase())
}

func (s *runnerActor) MakeBehaviour() Behaviour {
	return Behaviour{}
}

func NewRunnerActor(runner Runner) BehavioralActor {
	return &runnerActor{runner: runner}
}

type SimpleRunner func() error

type simpleRunnerActor struct {
	Actor
	runner SimpleRunner
}

func (s *simpleRunnerActor) Run() error {
	if s.runner == nil {
		return nil
	}
	return s.runner()
}

func (s *simpleRunnerActor) MakeBehaviour() Behaviour {
	return Behaviour{}
}

func NewSimpleRunnerActor(simpleRunner SimpleRunner) BehavioralActor {
	return &simpleRunnerActor{runner: simpleRunner}
}

type simpleNamedRunnerActor struct {
	simpleRunnerActor
	name string
}

func (s *simpleNamedRunnerActor) MakeBehaviour() Behaviour {
	return Behaviour{Name: s.name}
}

func NewSimpleNamedRunnerActor(name string, simpleRunner SimpleRunner) BehavioralActor {
	return &simpleNamedRunnerActor{simpleRunnerActor{runner: simpleRunner}, name}
}

type simpleCallbackRunnerActor struct {
	Actor
	runner common.SimpleCallback
}

func (s *simpleCallbackRunnerActor) Run() error {
	if s.runner == nil {
		return nil
	}
	s.runner()
	return nil
}

func (s *simpleCallbackRunnerActor) MakeBehaviour() Behaviour {
	return Behaviour{}
}

func NewSimpleCallbackRunnerActor(simpleCallback common.SimpleCallback) BehavioralActor {
	return &simpleCallbackRunnerActor{runner: simpleCallback}
}
