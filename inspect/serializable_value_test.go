package inspect_test

import (
	"gerrit-share.lan/go/inspect"
)

type Value string

func (v *Value) Inspect(inspector *inspect.GenericInspector) {
	inspector.String((*string)(v), "example_value", "this is example value type")
}

func ExampleInspectable_inspectValue() {
}
