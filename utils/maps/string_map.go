package maps

import ()

type String map[string]string

func (s *String) Add(key string, value string) {
	if *s == nil {
		*s = make(String, 1)
	}
	(*s)[key] = value
}

func (s *String) Update(key string, value string) bool {
	oldValue, ok := (*s)[key]
	if ok && oldValue == value {
		return false
	}
	s.Add(key, value)
	return true
}

func (s *String) Delete(key string) {
	delete(*s, key)
}

func (s *String) DeleteIfEquals(key string, value string) bool {
	oldValue, ok := (*s)[key]
	if ok && oldValue == value {
		delete(*s, key)
		return true
	}
	return false
}

func (s *String) Clear() {
	*s = nil
}

func (s String) IsEmpty() bool {
	return len(s) == 0
}

func (s String) Clone() String {
	clone := make(String)
	for key, value := range s {
		clone[key] = value
	}
	return clone
}

func (s *String) Append(other String) {
	for key, value := range other {
		s.Add(key, value)
	}
}
