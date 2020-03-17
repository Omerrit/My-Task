package tobytes

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/json/tojson"
)

type inspector struct {
	tojson.Inspector
	output     []byte
	inProgress bool
}

func (i *inspector) ValueString(value *string, typeName string, typeDescription string) error {
	if i.inProgress {
		return i.Inspector.ValueString(value, typeName, typeDescription)
	}
	if value == nil {
		return nil
	}
	i.output = []byte(*value)
	return nil
}

func (i *inspector) ValueByteString(value *[]byte, typeName string, typeDescription string) error {
	if i.inProgress {
		return i.Inspector.ValueByteString(value, typeName, typeDescription)
	}
	if value == nil {
		return nil
	}
	i.output = *value
	return nil
}

func (i *inspector) StartObject(typeName string, typeDescription string) error {
	i.inProgress = true
	return i.Inspector.StartObject(typeName, typeDescription)
}

func (i *inspector) StartArray(typeName string, valueTypeName string, typeDescription string) error {
	i.inProgress = true
	return i.Inspector.StartArray(typeName, valueTypeName, typeDescription)
}

func (i *inspector) StartMap(typeName string, valueTypeName string, typeDescription string) error {
	i.inProgress = true
	return i.Inspector.StartMap(typeName, valueTypeName, typeDescription)
}

func ToBytes(value inspect.Inspectable) ([]byte, error) {
	insp := &inspector{}
	serializer := inspect.NewGenericInspector(insp)
	value.Inspect(serializer)
	if len(insp.output) > 0 {
		return insp.output, serializer.GetError()
	}
	return insp.Inspector.Output(), serializer.GetError()
}
