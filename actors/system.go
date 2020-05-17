package actors

import (
	"gerrit-share.lan/go/actors/internal/types"
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/debug"
	//	"gerrit-share.lan/go/errors"
	"sync"
)

//TODO: make actor registry (actor->id, id->actor) with sync.Map
//to be used for remote messaging in the future

//need to make everything inspectable for this including all commands and messages
//and make special proxy actor with very limited functionality
//(can't generate or process commands except limited amount of basic ones, only forwarding)

type BehavioralActor interface {
	ActorCompatible
	MakeBehaviour() Behaviour
}

type System struct {
	serviceCounter    sync.WaitGroup
	types             types.Types
	surveillanceActor ActorService
	plugins           map[string]ActorService
}

//manual initialization, will be called automatically on the first Spawn() or Become().
//Call this manually if you plan to call first Spawns in parallel
func (s *System) Init() {
	types.Init(&s.types)
	s.startPlugins()
}

func (s *System) startPlugins() {
	if s.plugins != nil {
		return
	}
	s.plugins = make(map[string]ActorService)
	iteratePlugins(func(name string, creator ActorCreator) bool {
		if creator != nil {
			s.plugins[name] = s.makeService(creator())
		}
		return true
	})
	if debug.DebugLevelValue == debug.Debug {
		s.EnableSurveillance()
	}
	if s.surveillanceActor != nil {
		for _, service := range s.plugins {
			s.surveillanceActor.SendMessage(service)
		}
	}
	//everything must be initialized before first spawn
	for _, service := range s.plugins {
		go service.start(s)
	}
}

func (s *System) stopPlugins() {
	for _, actor := range s.plugins {
		actor.SendQuit(nil)
	}
}

func (s *System) makeService(actor BehavioralActor) ActorService {
	types.Init(&s.types) //make sure type system is initialized when the actor starts
	base := actor.GetBase()
	service := base.Service()
	service.init(actor, nil)
	base.init(s)
	s.serviceCounter.Add(1)
	if s.surveillanceActor != nil {
		s.surveillanceActor.SendMessage(service)
	}
	//	log.Printf("[%p] spawning from\n", service)
	//	log.Println(errors.CallerStack(3))
	return service
}

func (s *System) Spawn(actor BehavioralActor) ActorService {
	s.startPlugins()
	service := s.makeService(actor)
	go service.start(s)
	return service
}

func (s *System) Become(actor BehavioralActor) ActorService {
	s.startPlugins()
	service := s.makeService(actor)
	service.start(s)
	return service
}

//spawn function based actor
func (s *System) SpawnFunc(behaviourMaker BehaviourMaker) ActorService {
	return s.Spawn(NewSimpleActor(behaviourMaker))
}

//run function basic actor in current goroutine
func (s *System) BecomeFunc(behaviourMaker BehaviourMaker) ActorService {
	return s.Become(NewSimpleActor(behaviourMaker))
}

//spawn function based actor with empty behaviour
func (s *System) DoAsync(initializer SimpleInitializer) ActorService {
	return s.Spawn(NewSimpleInitializerActor(initializer))
}

//run function based actor in current goroutine with empty behaviour
func (s *System) Do(initializer SimpleInitializer) ActorService {
	return s.Become(NewSimpleInitializerActor(initializer))
}

//run function based namedactor in current goroutine with empty behaviour
func (s *System) DoNamed(name string, initializer SimpleInitializer) ActorService {
	return s.Become(NewSimpleNamedInitializerActor(name, initializer))
}

//spawn continuous function as actor, it should process actor messages if it wants to be responsive
func (s *System) RunAsync(runner Runner) ActorService {
	return s.Spawn(NewRunnerActor(runner))
}

//run continuous function as actor in current goroutine, it should process actor messages if it wants to be responsive
func (s *System) Run(runner Runner) ActorService {
	return s.Become(NewRunnerActor(runner))
}

//spawn continuous function that cannot process actor commands as actor,
//linking and monitoring would work that way but nothing else
func (s *System) RunAsyncSimple(simpleRunner SimpleRunner) ActorService {
	return s.Spawn(NewSimpleRunnerActor(simpleRunner))
}

//spawn continuous function that cannot process actor commands as named actor,
//linking and monitoring would work that way but nothing else
func (s *System) RunAsyncSimpleNamed(name string, simpleRunner SimpleRunner) ActorService {
	return s.Spawn(NewSimpleNamedRunnerActor(name, simpleRunner))
}

//spawn continuous function that cannot process actor commands as actor,
//linking and monitoring would work that way but nothing else
func (s *System) RunAsyncSimpleCallback(simpleCallback common.SimpleCallback) ActorService {
	return s.Spawn(NewSimpleCallbackRunnerActor(simpleCallback))
}

func (s *System) serviceFinished() {
	s.serviceCounter.Done()
}

func (s *System) WaitFinished() {
	if s.surveillanceActor != nil {
		s.surveillanceActor.SendMessage(closeMe{})
	}
	s.stopPlugins()
	s.serviceCounter.Wait()
}

func (s *System) Types() types.Types {
	return s.types
}

func (s *System) EnableSurveillance() {
	if s.surveillanceActor != nil {
		return
	}
	s.surveillanceActor = s.makeService(new(surveillanceActor))
	go s.surveillanceActor.start(s)
}

func (s *System) SubscribeToSurveillanceFeed(me ActorCompatible, input StreamInput, onError common.ErrorCallback) {
	me.GetBase().RequestStream(input, s.surveillanceActor, new(requestActorStream), onError)
}

func (s *System) GetPluginActor(name string) ActorService {
	return s.plugins[name]
}

//then the usual autostarter things
//then stream of actor services (input/output) to push published actors
//publishing should be done via plugin, not main system, only actor registry should be in the main system (and it should be THE main actor)
//monitoring system should also be done via plugins

//commands should be registered in system's type system to make them readable by type name
//types should also be registered globally (maybe)
