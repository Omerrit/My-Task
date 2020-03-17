package autorestarter

import ()

type ServiceMakers map[string]ServiceMaker

func (s *ServiceMakers) Add(name string, maker ServiceMaker) {
	if *s == nil {
		*s = make(ServiceMakers, 1)
	}
	(*s)[name] = maker
}

func (s *ServiceMakers) Remove(name string) {
	delete(*s, name)
}
