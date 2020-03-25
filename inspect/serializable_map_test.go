package inspect_test

import (
	"gerrit-share.lan/go/inspect"
)

type Map map[string]int

func (m *Map) setLength(length int) {
	if *m == nil {
		*m = make(Map, length)
		return
	}
	//map clearing idiom, supported by compiler
	for key := range *m {
		delete(*m, key)
	}
}

func (m *Map) Inspect(inspector *inspect.GenericInspector) {
	mi := inspector.Map("simple_map", "int", "simple integer map example")
	{
		if mi.IsReading() {
			length := mi.GetLength()
			m.setLength(length)
			for i := 0; i < length; i++ {
				(*m)[mi.NextKey()] = mi.ReadInt()
			}
		} else {
			mi.SetLength(len(*m))
			for key, value := range *m {
				mi.WriteInt(key, value)
			}
		}
		mi.End()
	}
}

//Object from 'inspect object' example
type ObjectMap map[string]Object

func (m *ObjectMap) setLength(length int) {
	if *m == nil {
		*m = make(ObjectMap, length)
		return
	}
	for key := range *m {
		delete(*m, key)
	}
}

func (m *ObjectMap) Inspect(inspector *inspect.GenericInspector) {
	mi := inspector.Map("object_map", ObjectName, "object map example")
	{
		if mi.IsReading() {
			length := mi.GetLength()
			m.setLength(length)
			for i := 0; i < length; i++ {
				key := mi.NextKey()
				var obj Object
				obj.Inspect(mi.ReadValue())
				(*m)[key] = obj
			}
		} else {
			mi.SetLength(len(*m))
			for key, value := range *m {
				value.Inspect(mi.WriteValue(key))
			}
		}
	}
}

func ExampleInspectable_inspectMap() {}
