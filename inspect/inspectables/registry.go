package inspectables

import (
	"gerrit-share.lan/go/inspect"
	"sync"
)

type Creator func() inspect.Inspectable

//map string->creator
//may be extended later
var registry sync.Map

type info struct {
	creator     Creator
	description string
}

func Register(name string, creator Creator) {
	if creator == nil {
		return
	}
	registry.Store(name, info{creator, ""})
}

func RegisterDescribed(name string, creator Creator, description string) {
	if creator == nil {
		return
	}
	registry.Store(name, info{creator, description})
}

func Get(name string) Creator {
	value, ok := registry.Load(name)
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	return value.(info).creator
}

func GetDescription(name string) string {
	value, ok := registry.Load(name)
	if !ok {
		return ""
	}
	if value == nil {
		return ""
	}
	return value.(info).description
}

//typed variant of sync.Map.Range
//parameters are name, creator function and description
func ForEach(f func(string, Creator, string) bool) {
	registry.Range(func(key interface{}, value interface{}) bool {
		i := value.(info)
		return f(key.(string), i.creator, i.description)
	})
}
