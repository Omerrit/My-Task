package tojson

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectorhelpers"
	"math/big"
)

type Inspector struct {
	inspectorhelpers.Writer
	jsonWriter
}

func (i *Inspector) ValueInt(value *int, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendInt(int64(*value))
	return nil
}

func (i *Inspector) ValueInt32(value *int32, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendInt(int64(*value))
	return nil
}

func (i *Inspector) ValueInt64(value *int64, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendInt(*value)
	return nil
}

func (i *Inspector) ValueFloat32(value *float32, format byte, precision int, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendFloat32(*value, format, precision)
	return nil
}

func (i *Inspector) ValueFloat64(value *float64, format byte, precision int, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendFloat64(*value, format, precision)
	return nil
}

func (i *Inspector) ValueString(value *string, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(*value)
	return nil
}

func (i *Inspector) ValueByteString(value *[]byte, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(string(*value))
	return nil
}

func (i *Inspector) ValueBytes(value *[]byte, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendBytes(*value)
	return nil
}

func (i *Inspector) ValueBool(value *bool, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendBool(*value)
	return nil
}

func (i *Inspector) ValueBigInt(value *big.Int, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	data, _ := value.MarshalText()
	i.appendString(string(data))
	return nil
}

func (i *Inspector) ValueRat(value *big.Rat, precision int, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(value.FloatString(precision))
	return nil
}

func (i *Inspector) ValueBigFloat(value *big.Float, format byte, precision int, typeName string, typeDescription string) error {
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(value.Text(format, precision))
	return nil
}

func (i *Inspector) StartObject(typeName string, typeDescription string) error {
	i.startObject()
	return nil
}

func (i *Inspector) StartArray(typeName string, valueTypeName string, typeDescription string) error {
	i.startArray()
	return nil
}

func (i *Inspector) StartMap(typeName string, valueTypeName string, typeDescription string) error {
	i.startObject()
	return nil
}

func (i *Inspector) EndObject() error {
	i.endObject()
	return nil
}

func (i *Inspector) EndArray() error {
	i.endArray()
	return nil
}

func (i *Inspector) EndMap() error {
	i.endObject()
	return nil
}

func (i *Inspector) ObjectInt(value *int, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendInt(int64(*value))
	return nil
}

func (i *Inspector) ObjectInt32(value *int32, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendInt(int64(*value))
	return nil
}

func (i *Inspector) ObjectInt64(value *int64, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendInt(*value)
	return nil
}

func (i *Inspector) ObjectFloat32(value *float32, format byte, precision int, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendFloat32(*value, format, precision)
	return nil
}

func (i *Inspector) ObjectFloat64(value *float64, format byte, precision int, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendFloat64(*value, format, precision)
	return nil
}

func (i *Inspector) ObjectString(value *string, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(*value)
	return nil
}

func (i *Inspector) ObjectByteString(value *[]byte, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(string(*value))
	return nil
}

func (i *Inspector) ObjectBytes(value *[]byte, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendBytes(*value)
	return nil
}

func (i *Inspector) ObjectBool(value *bool, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendBool(*value)
	return nil
}

func (i *Inspector) ObjectBigInt(value *big.Int, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	data, _ := value.MarshalText()
	i.appendString(string(data))
	return nil
}

func (i *Inspector) ObjectRat(value *big.Rat, precision int, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(value.FloatString(precision))
	return nil
}

func (i *Inspector) ObjectBigFloat(value *big.Float, format byte, precision int, name string, mandatory bool, description string) error {
	i.appendDelimiter()
	i.appendKey(name)
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(value.Text(format, precision))
	return nil
}

func (i *Inspector) ObjectValue(name string, mandatory bool, description string) (bool, error) {
	i.appendDelimiter()
	i.appendKey(name)
	return true, nil
}

func (i *Inspector) ArrayLen(length int) (int, error) {
	return length, nil
}

func (i *Inspector) ArrayInt(value *int) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendInt(int64(*value))
	return nil
}

func (i *Inspector) ArrayInt32(value *int32) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendInt(int64(*value))
	return nil
}

func (i *Inspector) ArrayInt64(value *int64) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendInt(*value)
	return nil
}

func (i *Inspector) ArrayFloat32(value *float32, format byte, precision int) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendFloat32(*value, format, precision)
	return nil
}

func (i *Inspector) ArrayFloat64(value *float64, format byte, precision int) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendFloat64(*value, format, precision)
	return nil
}

func (i *Inspector) ArrayString(value *string) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(*value)
	return nil
}

func (i *Inspector) ArrayByteString(value *[]byte) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(string(*value))
	return nil
}

func (i *Inspector) ArrayBytes(value *[]byte) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendBytes(*value)
	return nil
}

func (i *Inspector) ArrayBool(value *bool) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendBool(*value)
	return nil
}

func (i *Inspector) ArrayBigInt(value *big.Int) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	data, _ := value.MarshalText()
	i.appendString(string(data))
	return nil
}

func (i *Inspector) ArrayRat(value *big.Rat, precision int) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(value.FloatString(precision))
	return nil
}

func (i *Inspector) ArrayBigFloat(value *big.Float, format byte, precision int) error {
	i.appendDelimiter()
	if value == nil {
		i.appendNull()
		return nil
	}
	i.appendString(value.Text(format, precision))
	return nil
}

func (i *Inspector) ArrayValue() error {
	i.appendDelimiter()
	return nil
}

func (i *Inspector) MapLen(length int) (int, error) {
	return length, nil
}

func (i *Inspector) MapWriteInt(key string, value int) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendInt(int64(value))
	return nil
}

func (i *Inspector) MapWriteInt32(key string, value int32) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendInt(int64(value))
	return nil
}

func (i *Inspector) MapWriteInt64(key string, value int64) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendInt(value)
	return nil
}

func (i *Inspector) MapWriteFloat32(key string, value float32, format byte, precision int) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendFloat32(value, format, precision)
	return nil
}

func (i *Inspector) MapWriteFloat64(key string, value float64, format byte, precision int) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendFloat64(value, format, precision)
	return nil
}

func (i *Inspector) MapWriteString(key string, value string) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendString(value)
	return nil
}

func (i *Inspector) MapWriteByteString(key string, value []byte) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendString(string(value))
	return nil
}

func (i *Inspector) MapWriteBytes(key string, value []byte) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendBytes(value)
	return nil
}

func (i *Inspector) MapWriteBool(key string, value bool) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendBool(value)
	return nil
}

func (i *Inspector) MapWriteBigInt(key string, value *big.Int) error {
	i.appendDelimiter()
	i.appendKey(key)
	data, _ := value.MarshalText()
	i.appendString(string(data))
	return nil
}

func (i *Inspector) MapWriteRat(key string, value *big.Rat, precision int) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendString(value.FloatString(precision))
	return nil
}

func (i *Inspector) MapWriteBigFloat(key string, value *big.Float, format byte, precision int) error {
	i.appendDelimiter()
	i.appendKey(key)
	i.appendString(value.Text(format, precision))
	return nil
}

func (i *Inspector) MapWriteValue(key string) error {
	i.appendDelimiter()
	i.appendKey(key)
	return nil
}

func init() {
	var _ inspect.InspectorImpl = (*Inspector)(nil)
}
