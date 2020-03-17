package multisets

import ()

type String map[string]int

func (s *String) Add(value string) {
	if *s == nil {
		*s = make(String)
	}
	(*s)[value] = (*s)[value] + 1
}

func (s *String) Remove(value string) {
	count, ok := (*s)[value]
	if ok {
		if count <= 1 {
			delete(*s, value)
		} else {
			(*s)[value] = count - 1
		}
	}
}

func (s *String) RemoveAll(value string) {
	delete(*s, value)
}

func (s *String) Contains(value string) bool {
	_, ok := (*s)[value]
	return ok
}

func (s *String) Count(value string) int {
	return (*s)[value]
}

func (s *String) Clear() {
	*s = make(String)
}

func (s String) IsEmpty() bool {
	return len(s) == 0
}
