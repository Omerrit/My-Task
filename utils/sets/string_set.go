package sets

import (
	"gerrit-share.lan/go/common"
)

type String map[string]common.None

func (s *String) Add(value string) {
	if *s == nil {
		*s = make(String)
	}
	(*s)[value] = common.None{}
}

func (s *String) Remove(value string) {
	delete(*s, value)
}

func (s *String) Contains(value string) bool {
	_, ok := (*s)[value]
	return ok
}

func (s *String) Clear() {
	*s = make(String)
}

func (s String) IsEmpty() bool {
	return len(s) == 0
}

func (s *String) Union(other String) {
	for value, _ := range other {
		s.Add(value)
	}
}

func Join(first String, second String) String {
	result := String{}
	for value, _ := range first {
		result.Add(value)
	}
	for value, _ := range second {
		result.Add(value)
	}
	return result
}
