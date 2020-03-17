package interfaces

import (
	"gerrit-share.lan/go/common"
)

type Named interface {
	Name() string
}

type DummyNamed common.None

func (DummyNamed) Name() string {
	return "dummy"
}

type simpleNamed struct {
	name string
}

func (s simpleNamed) Name() string {
	return s.name
}

func MakeSimpleNamed(name string) Named {
	return simpleNamed{name}
}

func GetName(object interface{}) string {
	var result string
	named, ok := object.(Named)
	if ok {
		result = named.Name()
	}
	return result
}
