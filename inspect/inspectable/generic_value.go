package inspectable

import (
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"math/big"
)

type GenericValue struct {
	value  interface{}
	typeId inspect.TypeId
}

func NewGenericValue(typeId inspect.TypeId) *GenericValue {
	return &GenericValue{typeId: typeId}
}

func (g *GenericValue) SetValue(value interface{}) {
	g.value = value
}

func (g *GenericValue) Inspect(i *inspect.GenericInspector) {
	if i.IsReading() {
		return
	}
	switch g.typeId {
	case inspect.TypeMap, inspect.TypeObject, inspect.TypeArray, inspect.TypeValue:
		v, ok := g.value.(inspect.Inspectable)
		if !ok {
			i.SetError(errors.Errorf("generic value: not inspectable value detected for type id: %#v", g.typeId))
			return
		}
		v.Inspect(i)
	case inspect.TypeString:
		v := (g.value).(string)
		i.String(&v, "", "")
	case inspect.TypeInt:
		v := (g.value).(int)
		i.Int(&v, "", "")
	case inspect.TypeInt32:
		v := (g.value).(int32)
		i.Int32(&v, "", "")
	case inspect.TypeInt64:
		v := (g.value).(int64)
		i.Int64(&v, "", "")
	case inspect.TypeFloat32:
		v := (g.value).(float32)
		i.Float32(&v, 'g', -1, "", "")
	case inspect.TypeFloat64:
		v := (g.value).(float64)
		i.Float64(&v, 'g', -1, "", "")
	case inspect.TypeBool:
		v := (g.value).(bool)
		i.Bool(&v, "", "")
	case inspect.TypeByteString:
		v := (g.value).([]byte)
		i.ByteString(&v, "", "")
	case inspect.TypeBytes:
		v := (g.value).([]byte)
		i.Bytes(&v, "", "")
	case inspect.TypeRat:
		v := (g.value).(*big.Rat)
		i.Rat(v, -1, "", "")
	case inspect.TypeBigInt:
		v := (g.value).(*big.Int)
		i.BigInt(v, "", "")
	case inspect.TypeBigFloat:
		v := (g.value).(*big.Float)
		i.BigFloat(v, 'g', -1, "", "")
	case inspect.TypeInvalid:
		v := "ok"
		i.String(&v, "", "")
	default:
		panic("unexpected type id")
	}
}
