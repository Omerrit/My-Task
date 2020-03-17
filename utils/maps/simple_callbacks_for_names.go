package maps

import (
	"gerrit-share.lan/go/common"
)

type SimpleCallbacksForNames map[string]common.SimpleCallback

func (s *SimpleCallbacksForNames) Add(name string, callback common.SimpleCallback) {
	if *s == nil {
		*s = make(SimpleCallbacksForNames, 1)
	}
	(*s)[name] = callback
}

func (s *SimpleCallbacksForNames) Delete(name string) {
	delete(*s, name)
}

func (s *SimpleCallbacksForNames) Clear() {
	*s = nil
}

func (s SimpleCallbacksForNames) IsEmpty() bool {
	return len(s) == 0
}
