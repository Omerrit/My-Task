package sets

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/utils/sets"
)

type StringsForServices map[actors.ActorService]sets.String

func (s *StringsForServices) Add(service actors.ActorService, value string) {
	if *s == nil {
		*s = make(StringsForServices, 1)
	}
	set := (*s)[service]
	set.Add(value)
	(*s)[service] = set
}

func (s *StringsForServices) Remove(service actors.ActorService, value string) {
	set := (*s)[service]
	set.Remove(value)
	if set.IsEmpty() {
		s.RemoveByDestination(service)
	} else {
		(*s)[service] = set
	}
}

func (s *StringsForServices) RemoveByDestination(destination actors.ActorService) {
	delete(*s, destination)
}

func (s *StringsForServices) Clear() {
	*s = nil
}
