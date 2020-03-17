package inspectwrappers

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/mreflect"
	"math/big"
)

const packageName = "inspect.wrappers"

const defaultFormat byte = 'g'
const defaultPrecision int = 16

type Value struct {
	typeId     inspect.TypeId
	realTypeId mreflect.TypeId
	isPointer  bool
	data       interface{}
	format     byte
	precision  int
}

const ValueName = packageName + ".value"

func NewIntValue(value int) *Value {
	return &Value{inspect.TypeInt, mreflect.GetTypeId(value), false, value, 0, 0}
}

func NewPIntValue(value *int) *Value {
	return &Value{inspect.TypeInt, mreflect.GetTypeId(value), true, value, 0, 0}
}

func NewInt32Value(value int32) *Value {
	return &Value{inspect.TypeInt32, mreflect.GetTypeId(value), false, value, 0, 0}
}

func NewPInt32Value(value *int32) *Value {
	return &Value{inspect.TypeInt32, mreflect.GetTypeId(value), true, value, 0, 0}
}

func NewInt64Value(value int64) *Value {
	return &Value{inspect.TypeInt64, mreflect.GetTypeId(value), false, value, 0, 0}
}

func NewPInt64Value(value *int64) *Value {
	return &Value{inspect.TypeInt64, mreflect.GetTypeId(value), true, value, 0, 0}
}

func NewFloat32Value(value float32, format byte, precision int) *Value {
	return &Value{inspect.TypeFloat32, mreflect.GetTypeId(value), false, value, format, precision}
}

func NewPFloat32Value(value *float32, format byte, precision int) *Value {
	return &Value{inspect.TypeFloat32, mreflect.GetTypeId(value), true, value, format, precision}
}

func NewFloat64Value(value float64, format byte, precision int) *Value {
	return &Value{inspect.TypeFloat32, mreflect.GetTypeId(value), false, value, format, precision}
}

func NewPFloat64Value(value *float64, format byte, precision int) *Value {
	return &Value{inspect.TypeFloat32, mreflect.GetTypeId(value), true, value, format, precision}
}

func NewStringValue(value string) *Value {
	return &Value{inspect.TypeString, mreflect.GetTypeId(value), false, value, 0, 0}
}

func NewPStringValue(value *string) *Value {
	return &Value{inspect.TypeString, mreflect.GetTypeId(value), true, value, 0, 0}
}

func NewByteStringValue(value []byte) *Value {
	return &Value{inspect.TypeByteString, mreflect.GetTypeId(value), false, value, 0, 0}
}

func NewPByteStringValue(value *[]byte) *Value {
	return &Value{inspect.TypeByteString, mreflect.GetTypeId(value), true, value, 0, 0}
}

func NewBytesValue(value []byte) *Value {
	return &Value{inspect.TypeBytes, mreflect.GetTypeId(value), false, value, 0, 0}
}

func NewPBytesValue(value *[]byte) *Value {
	return &Value{inspect.TypeBytes, mreflect.GetTypeId(value), true, value, 0, 0}
}

func NewBoolValue(value bool) *Value {
	return &Value{inspect.TypeBool, mreflect.GetTypeId(value), false, value, 0, 0}
}

func NewPBoolValue(value *bool) *Value {
	return &Value{inspect.TypeBool, mreflect.GetTypeId(value), true, value, 0, 0}
}

func NewBigIntValue(value *big.Int) *Value {
	return &Value{inspect.TypeBigInt, mreflect.GetTypeId(value), true, value, 0, 0}
}

func NewRatValue(value *big.Rat, precision int) *Value {
	return &Value{inspect.TypeRat, mreflect.GetTypeId(value), true, value, 0, precision}
}

func NewBigFloatValue(value *big.Float, format byte, precision int) *Value {
	return &Value{inspect.TypeBigFloat, mreflect.GetTypeId(value), true, value, format, precision}
}

func NewInspectableValue(value inspect.Inspectable) *Value {
	return &Value{inspect.TypeValue, mreflect.GetTypeId(value), true, value, 0, 0}
}

//will panic if type is not supported
//doesn't support []byte due to ambiguity
func NewGuessValue(value interface{}, format byte, precision int) (*Value, error) {
	switch data := value.(type) {
	case int:
		return NewIntValue(data), nil
	case *int:
		return NewPIntValue(data), nil
	case int32:
		return NewInt32Value(data), nil
	case *int32:
		return NewPInt32Value(data), nil
	case int64:
		return NewInt64Value(data), nil
	case *int64:
		return NewPInt64Value(data), nil
	case float32:
		return NewFloat32Value(data, format, precision), nil
	case *float32:
		return NewPFloat32Value(data, format, precision), nil
	case float64:
		return NewFloat64Value(data, format, precision), nil
	case *float64:
		return NewPFloat64Value(data, format, precision), nil
	case string:
		return NewStringValue(data), nil
	case *string:
		return NewPStringValue(data), nil
	case bool:
		return NewBoolValue(data), nil
	case *bool:
		return NewPBoolValue(data), nil
	case *big.Int:
		return NewBigIntValue(data), nil
	case *big.Rat:
		return NewRatValue(data, precision), nil
	case *big.Float:
		return NewBigFloatValue(data, format, precision), nil
	case inspect.Inspectable:
		return NewInspectableValue(data), nil
	case []byte:
		return nil, ErrAmbiguity
	case *[]byte:
		return nil, ErrAmbiguity
	default:
		return nil, ErrUnsupported
	}
}

// create default value of a given type using go defaults
//doesn't support TypeValue
//resulting value can't be written to
func NewDefaultBasicValue(valueType inspect.TypeId) (*Value, error) {
	switch valueType {
	case inspect.TypeInt:
		return NewIntValue(0), nil
	case inspect.TypeInt32:
		return NewInt32Value(0), nil
	case inspect.TypeInt64:
		return NewInt64Value(0), nil
	case inspect.TypeFloat32:
		return NewFloat32Value(0, defaultFormat, defaultPrecision), nil
	case inspect.TypeFloat64:
		return NewFloat64Value(0, defaultFormat, defaultPrecision), nil
	case inspect.TypeString:
		return NewStringValue(""), nil
	case inspect.TypeByteString:
		return NewByteStringValue(nil), nil
	case inspect.TypeBytes:
		return NewBytesValue(nil), nil
	case inspect.TypeBool:
		return NewBoolValue(false), nil
	case inspect.TypeBigInt:
		return NewBigIntValue(new(big.Int)), nil
	case inspect.TypeRat:
		return NewRatValue(new(big.Rat), defaultPrecision), nil
	case inspect.TypeBigFloat:
		return NewBigFloatValue(new(big.Float), defaultFormat, defaultPrecision), nil
	default:
		return nil, ErrUnsupported
	}
}

func (v *Value) Replace(value interface{}) error {
	if mreflect.GetTypeId(value) != v.realTypeId {
		return ErrWrongType
	}
	v.data = value
	return nil
}

func (v *Value) SetFormat(format byte) {
	v.format = format
}

func (v *Value) SetPrecision(precision int) {
	v.precision = precision
}

func (v *Value) SetFormatPrecision(format byte, precision int) {
	v.format = format
	v.precision = precision
}

func (v *Value) Inspect(inspector *inspect.GenericInspector) {
	if inspector.IsReading() && !v.isPointer {
		inspector.SetError(ErrCantWrite)
		return
	}
	switch v.typeId {
	case inspect.TypeInt:
		if v.isPointer {
			inspector.Int(v.data.(*int), ValueName, "")
		} else {
			i := v.data.(int)
			inspector.Int(&i, ValueName, "")
		}
	case inspect.TypeInt32:
		if v.isPointer {
			inspector.Int32(v.data.(*int32), ValueName, "")
		} else {
			i := v.data.(int32)
			inspector.Int32(&i, ValueName, "")
		}
	case inspect.TypeInt64:
		if v.isPointer {
			inspector.Int64(v.data.(*int64), ValueName, "")
		} else {
			i := v.data.(int64)
			inspector.Int64(&i, ValueName, "")
		}
	case inspect.TypeFloat32:
		if v.isPointer {
			inspector.Float32(v.data.(*float32), v.format, v.precision, ValueName, "")
		} else {
			i := v.data.(float32)
			inspector.Float32(&i, v.format, v.precision, ValueName, "")
		}
	case inspect.TypeFloat64:
		if v.isPointer {
			inspector.Float64(v.data.(*float64), v.format, v.precision, ValueName, "")
		} else {
			i := v.data.(float64)
			inspector.Float64(&i, v.format, v.precision, ValueName, "")
		}
	case inspect.TypeString:
		if v.isPointer {
			inspector.String(v.data.(*string), ValueName, "")
		} else {
			i := v.data.(string)
			inspector.String(&i, ValueName, "")
		}
	case inspect.TypeByteString:
		if v.isPointer {
			inspector.ByteString(v.data.(*[]byte), ValueName, "")
		} else {
			i := v.data.([]byte)
			inspector.ByteString(&i, ValueName, "")
		}
	case inspect.TypeBytes:
		if v.isPointer {
			inspector.Bytes(v.data.(*[]byte), ValueName, "")
		} else {
			i := v.data.([]byte)
			inspector.Bytes(&i, ValueName, "")
		}
	case inspect.TypeBool:
		if v.isPointer {
			inspector.Bool(v.data.(*bool), ValueName, "")
		} else {
			i := v.data.(bool)
			inspector.Bool(&i, ValueName, "")
		}
	case inspect.TypeBigInt:
		inspector.BigInt(v.data.(*big.Int), ValueName, "")
	case inspect.TypeRat:
		inspector.Rat(v.data.(*big.Rat), v.precision, ValueName, "")
	case inspect.TypeBigFloat:
		inspector.BigFloat(v.data.(*big.Float), v.format, v.precision, ValueName, "")
	case inspect.TypeValue:
		v.data.(inspect.Inspectable).Inspect(inspector)
	}
}

func init() {
	inspect.TestInspectable((*Value)(nil))
}
