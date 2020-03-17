package sets

import (
	"gerrit-share.lan/go/common"
)

type Int map[int]common.None

func (i *Int) Add(value int) {
	if *i == nil {
		*i = make(Int)
	}
	(*i)[value] = common.None{}
}

func (i *Int) Remove(value int) {
	delete(*i, value)
}

func (i *Int) Contains(value int) bool {
	_, ok := (*i)[value]
	return ok
}

func (i *Int) Clear() {
	*i = nil
}

func (i Int) IsEmpty() bool {
	return len(i) == 0
}
