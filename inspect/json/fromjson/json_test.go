package fromjson

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/json/tojson"
	"testing"
)

type fpvalue float32

func (v *fpvalue) Inspect(i *inspect.GenericInspector) {
	i.Float32((*float32)(v), 'g', -1, "fpvalue", "")
}

type simple struct {
	name    string
	value   int
	fpvalue fpvalue
}

func (s *simple) Inspect(inspector *inspect.GenericInspector) {
	i := inspector.Object("simple", "")
	i.String(&s.name, "name", true, "")
	i.Int(&s.value, "value", true, "")
	s.fpvalue.Inspect(i.Value("fpvalue", true, ""))
	i.End()
}

var testObject = simple{`he\"llo`, 3, 1.34}

func TestObjectDecode(t *testing.T) {
	var serializer tojson.Inspector
	testObject.Inspect(inspect.NewGenericInspector(&serializer))
	var value simple
	t.Log(string(serializer.Output()))
	reader := inspect.NewGenericInspector(NewInspector(serializer.Output(), 0))
	value.Inspect(reader)
	t.Log(value.name)
	if value != testObject {
		t.Errorf("read: %#v, expected: %#v, error:%v", value, testObject, reader.GetError())
	}
	t.Log(value)
}

type testArray []simple

func (t *testArray) Inspect(inspector *inspect.GenericInspector) {
	i := inspector.Array("testArray", "simple", "")
	if i.IsReading() {
		*t = make(testArray, i.GetLength())
	} else {
		i.SetLength(len(*t))
	}
	for index := range *t {
		(*t)[index].Inspect(i.Value())
	}
	i.End()
}

var arrayForTest = testArray{{`hello`, 3, 1.34}, {`hi`, 4, 1.35}, {"\t\t\"", 5, 1.36}}

func TestArrayDecode(t *testing.T) {
	var serializer tojson.Inspector
	arrayForTest.Inspect(inspect.NewGenericInspector(&serializer))
	var result testArray
	t.Log(string(serializer.Output()))
	reader := inspect.NewGenericInspector(NewInspector(serializer.Output(), 0))
	result.Inspect(reader)
	if len(result) == len(arrayForTest) {
		for i := 0; i < len(result); i++ {
			if result[i] != arrayForTest[i] {
				t.Errorf("index: %d, read: %#v, expected: %#v", i, result[i], arrayForTest[i])
			}
		}
	} else {
		t.Errorf("Length differ, have: %d, should be %d", len(result), len(arrayForTest))
	}
	if reader.GetError() != nil {
		t.Errorf("%#v,%s", reader.GetError(), reader.GetError())
	}
}

type objmap map[string]simple

func (o *objmap) Inspect(inspector *inspect.GenericInspector) {
	i := inspector.Map("objmap", "simple", "")
	if i.IsReading() {
		*o = make(objmap, i.GetLength())
		for index := 0; index < i.GetLength(); index++ {
			key := i.NextKey()
			var value simple
			value.Inspect(i.ReadValue())
			(*o)[key] = value
		}
	} else {
		i.SetLength(len(*o))
		for key, value := range *o {
			value.Inspect(i.WriteValue(key))
		}
	}
	i.End()
}

func (o objmap) IsSame(other objmap) bool {
	if len(o) != len(other) {
		return false
	}
	for key, value := range o {
		if other[key] != value {
			return false
		}
	}
	return true
}

var testMap = objmap{"one": simple{`hello`, 3, 1.34}, "two": simple{`hi`, 4, 1.35}, "three": simple{"\t\t\"", 5, 1.36}}

func TestMapDecode(t *testing.T) {
	var serializer tojson.Inspector
	testMap.Inspect(inspect.NewGenericInspector(&serializer))
	var result objmap
	t.Log(string(serializer.Output()))
	reader := inspect.NewGenericInspector(NewInspector(serializer.Output(), 0))
	result.Inspect(reader)
	if reader.GetError() != nil {
		t.Errorf("%#v,%s", reader.GetError(), reader.GetError())
	}
	if !result.IsSame(testMap) {
		t.Errorf("results differ: read: %#v, expected: %#v", result, testMap)
	}
}
