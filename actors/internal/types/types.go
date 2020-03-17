package types

import (
	"gerrit-share.lan/go/inspect/inspectables"
	"gerrit-share.lan/go/mreflect"
)

//inspectable may have Fill() method or may not
//these structs should be filed from inspectables database
type typeInfo struct {
	name        string
	creator     inspectables.Creator
	description string
}

type Types struct {
	//to be used by deserializer
	byName map[string]typeInfo
	//to be used by behaviour compiler (it works with raw types)
	byId map[mreflect.TypeId]typeInfo
}

func Init(t *Types) {
	if t.byId != nil || t.byName != nil {
		return
	}
	t.byName = make(map[string]typeInfo)
	t.byId = make(map[mreflect.TypeId]typeInfo)
	inspectables.ForEach(func(name string, creator inspectables.Creator, description string) bool {
		t.byName[name] = typeInfo{name, creator, description}
		t.byId[mreflect.GetTypeId(creator())] = typeInfo{name, creator, description}
		return true
	})
}

func (t *Types) IsRegisteredId(id mreflect.TypeId) bool {
	_, ok := t.byId[id]
	return ok
}

func (t *Types) IsRegistered(sample interface{}) bool {
	_, ok := t.byId[mreflect.GetTypeId(sample)]
	return ok
}

func (t *Types) GetNameById(id mreflect.TypeId) (string, bool) {
	info, ok := t.byId[id]
	return info.name, ok
}

func (t *Types) GetName(sample interface{}) (string, bool) {
	info, ok := t.byId[mreflect.GetTypeId(sample)]
	return info.name, ok
}

func (t *Types) GetNameDescription(name string) (string, bool) {
	info, ok := t.byName[name]
	return info.description, ok
}
