package sets

import (
	"gerrit-share.lan/go/common"
)

type Int64 map[int64]common.None

func (s *Int64) Add(value int64) {
	if *s == nil {
		*s = make(Int64)
	}
	(*s)[value] = common.None{}
}

func (s *Int64) Remove(value int64) {
	delete(*s, value)
}

func (s *Int64) Contains(value int64) bool {
	_, ok := (*s)[value]
	return ok
}

func (s *Int64) Clear() {
	*s = make(Int64)
}

func (s Int64) IsEmpty() bool {
	return len(s) == 0
}
