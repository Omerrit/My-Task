package tojson

import (
	"gerrit-share.lan/go/inspect"
	"testing"
)

type strarray []string

func (s *strarray) Inspect(inspector *inspect.GenericInspector) {
	a := inspector.Array("strarray", "string", "")
	a.SetLength(len(*s))
	for _, v := range *s {
		a.String(&v)
	}
	a.End()
}

func CompareStrings(t *testing.T, expected string, got string) {
	if expected != got {
		t.Error("expected", expected, "got", got)
	} else {
		t.Log(got)
	}
}

type StringBytesOutputter interface {
	Output() []byte
}

func CheckOutput(t *testing.T, expected string, outputter StringBytesOutputter) {
	CompareStrings(t, expected, string(outputter.Output()))
}

func TestArrayInspect(t *testing.T) {
	array := strarray{"hello", "world", "me"}
	expected := `["hello","world","me"]`
	var inspector Inspector
	array.Inspect(inspect.NewGenericInspector(&inspector))
	CheckOutput(t, expected, &inspector)
}

type fpvalue float32

func (fpv *fpvalue) Inspect(inspector *inspect.GenericInspector) {
	inspector.Float32((*float32)(fpv), 'g', -1, "fpvalue", "")
}

type object struct {
	id     int64
	value  fpvalue
	values strarray
}

func (o *object) Inspect(inspector *inspect.GenericInspector) {
	i := inspector.Object("object", "")
	i.Int64(&o.id, "id", true, "")
	o.value.Inspect(i.Value("value", true, ""))
	o.values.Inspect(i.Value("array", false, ""))
	i.End()
}

func TestObjectInspect(t *testing.T) {
	obj := object{10, 1.34, strarray{"hello", "world", "me"}}
	expected := `{"id":10,"value":1.34,"array":["hello","world","me"]}`
	var inspector Inspector
	obj.Inspect(inspect.NewGenericInspector(&inspector))
	CheckOutput(t, expected, &inspector)
}

type objarray []object

func (o *objarray) Inspect(inspector *inspect.GenericInspector) {
	i := inspector.Array("objarray", "object", "")
	i.SetLength(len(*o))
	for _, obj := range *o {
		obj.Inspect(i.Value())
	}
	i.End()
}

func TestObjectArrayInspect(t *testing.T) {
	arr := objarray{{10, 1.34, nil}, {11, 1.35, strarray{"hello"}}, {12, 1.36, strarray{"hello", "world"}}}
	expected := `[{"id":10,"value":1.34,"array":[]},{"id":11,"value":1.35,"array":["hello"]},{"id":12,"value":1.36,"array":["hello","world"]}]`
	var inspector Inspector
	arr.Inspect(inspect.NewGenericInspector(&inspector))
	CheckOutput(t, expected, &inspector)
}

type mapelem struct {
	key   string
	value object
}

//can't guarantee ordering for map type, hard to check serialized result
type objmap []mapelem

func (o *objmap) Inspect(inspector *inspect.GenericInspector) {
	i := inspector.Map("objmap", "object", "")
	i.SetLength(len(*o))
	for _, elem := range *o {
		elem.value.Inspect(i.WriteValue(elem.key))
	}
	i.End()
}

func TestObjectMapInspect(t *testing.T) {
	omap := objmap{{"one", object{10, 1.34, nil}}, {"two", object{11, 1.35, strarray{"hello"}}}, {"three", object{12, 1.36, strarray{"hello", "world"}}}}
	expected := `{"one":{"id":10,"value":1.34,"array":[]},"two":{"id":11,"value":1.35,"array":["hello"]},"three":{"id":12,"value":1.36,"array":["hello","world"]}}`
	var inspector Inspector
	omap.Inspect(inspect.NewGenericInspector(&inspector))
	CheckOutput(t, expected, &inspector)

}
