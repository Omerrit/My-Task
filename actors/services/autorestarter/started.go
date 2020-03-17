package autorestarter

import (
	"gerrit-share.lan/go/actors"
)

type started map[actors.ActorService]string

func (s *started) add(service actors.ActorService, name string) {
	if *s == nil {
		*s = make(started, 1)
	}
	(*s)[service] = name
}

func (s *started) remove(service actors.ActorService) {
	delete(*s, service)
}

func (s *started) contains(service actors.ActorService) bool {
	_, ok := (*s)[service]
	return ok
}
