package autorestarter

import (
	"errors"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/utils/sets"
	"log"
	"time"
)

type ServiceMaker func(parent *actors.Actor, name string) (actors.ActorService, error)

type staticAutorestarter struct {
	actors.Actor
	makers            ServiceMakers
	notStarted        sets.String
	started           started
	autorestartPeriod time.Duration
	name              string
}

func (s *staticAutorestarter) Shutdown() error {
	log.Println(s.name, "shut down")
	return nil
}

func (s *staticAutorestarter) scheduleRestart() {
	if !s.notStarted.IsEmpty() {
		time.AfterFunc(s.autorestartPeriod, func() { s.SendMessage(s.Service(), autostartMessage{}) })
	}
}

func (s *staticAutorestarter) Monitor(service actors.ActorService) {
	s.Actor.Monitor(service, func(err error) {
		if !s.started.contains(service) {
			return
		}
		if errors.Is(err, actors.ErrNotGonnaHappen) {
			if name, ok := s.started[service]; ok {
				log.Println("Service ", name, "stopped and doesn't want to be relaunched : ", err,
					"\nremoving from autorestart list")
				s.makers.Remove(name)
			}
			s.started.remove(service)
			return
		}
		s.notStarted.Add(s.started[service])
		s.started.remove(service)
		s.scheduleRestart()
	})
}

func (s *staticAutorestarter) autostart() {
	for name := range s.notStarted {
		service, err := s.makers[name](s.GetBase(), name)
		if err != nil {
			if errors.Is(err, actors.ErrNotGonnaHappen) {
				log.Println("Service ", name, "refused to launch : ", err,
					"\nremoving from autorestart list")
				s.makers.Remove(name)
				s.notStarted.Remove(name)
			} else {
				log.Println("Failed to start service ", name, " :", err)
			}
			continue
		}
		s.Monitor(service)
		s.started.add(service, name)
		s.notStarted.Remove(name)
	}
	s.scheduleRestart()
}

func (s *staticAutorestarter) MakeBehaviour() actors.Behaviour {
	log.Println(s.name, "started")
	for name, maker := range s.makers {
		service, err := maker(s.GetBase(), name)
		if err != nil {
			if errors.Is(err, actors.ErrNotGonnaHappen) {
				log.Println("Service ", name, "refused to launch : ", err,
					"\nremoving from autorestart list")
				s.makers.Remove(name)
			} else {
				log.Println("Failed to start service ", name, " :", err)
				s.notStarted.Add(name)
			}
			continue
		}
		s.Monitor(service)
		s.started.add(service, name)
		s.scheduleRestart()
	}
	b := actors.Behaviour{Name: s.name}
	b.AddMessage(autostartMessage{}, func(msg interface{}) {
		s.autostart()
	})
	return b
}

func NewStaticAutorestarter(s *actors.System, serviceName string, makers ServiceMakers, autorestartPeriod time.Duration) actors.ActorService {
	return s.Spawn(&staticAutorestarter{makers: makers, autorestartPeriod: autorestartPeriod, name: serviceName})
}
