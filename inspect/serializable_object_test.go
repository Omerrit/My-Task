package inspect_test

import (
	"gerrit-share.lan/go/inspect"
)

type Object struct {
	name  string
	value int
}

const ObjectName = "example_type" //to be used by containers for element type name

func (e *Object) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(ObjectName, "this is example type")
	{
		o.String(&e.name, "name", true, "name field")
		o.Int(&e.value, "value", true, "value field")
		o.End()
	}
}

func ExampleInspectable_inspectObject() {
}
