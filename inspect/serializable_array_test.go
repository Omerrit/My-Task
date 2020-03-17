package inspect_test

import (
	"gerrit-share.lan/go/inspect"
)

type Array []int

func (a *Array) setLength(length int) {
	if cap(*a) < length {
		*a = make(Array, length)
		return
	}
	*a = (*a)[:length]
}

func (a *Array) Inspect(inspector *inspect.GenericInspector) {
	ai := inspector.Array("simple_array", "int", "simple array example")
	{
		if ai.IsReading() {
			a.setLength(ai.GetLength())
		} else {
			ai.SetLength(len(*a))
		}
		for i := 0; i < len(*a); i++ {
			ai.Int(&(*a)[i])
		}
		ai.End()
	}
}

//Object from 'inspect object' example
type ObjectArray []Object

func (o *ObjectArray) setLength(length int) {
	if cap(*a) < length {
		*a = make(ObjectArray, length)
		return
	}
	*a = (*a)[:length]
}

func (o *ObjectArray) Inspect(inspector *inspect.GenericInspector) {
	ai := inspector.Array("object_array", ObjectName, "object array example")
	{
		if ai.IsReading() {
			o.setLength(ai.GetLength())
		} else {
			ai.SetLength(len(*o))
		}
		for i := 0; i < len(*o); i++ {
			(*o)[i].Inspect(ai.Value())
		}
		ai.End()
	}
}

func ExampleInspectable_inspectArray() {}
