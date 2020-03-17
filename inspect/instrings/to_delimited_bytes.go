package instrings

import (
	"encoding/base64"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/utils/simpleutils"
	"math/big"
	"strconv"
)

type ToDelimitedBytesImpl struct {
	Delimiter   []byte
	output      []byte
	hadAnything bool
}

func (t *ToDelimitedBytesImpl) IsReading() bool {
	return false
}

func (t *ToDelimitedBytesImpl) append(str []byte) {
	t.output = append(t.output, str...)
}

func (t *ToDelimitedBytesImpl) appendDelimiter() {
	if t.hadAnything {
		t.output = append(t.output, '\n') //t.Delimiter...)
	} else {
		t.hadAnything = true
	}
}

func (t *ToDelimitedBytesImpl) ArrayInt(value *int) error {
	t.appendDelimiter()
	if value != nil {
		t.output = strconv.AppendInt(t.output, int64(*value), 10)
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayInt32(value *int32) error {
	t.appendDelimiter()
	if value != nil {
		t.output = strconv.AppendInt(t.output, int64(*value), 10)
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayInt64(value *int64) error {
	t.appendDelimiter()
	if value != nil {
		t.output = strconv.AppendInt(t.output, *value, 10)
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayFloat32(value *float32, format byte, precision int) error {
	t.appendDelimiter()
	if value != nil {
		t.output = strconv.AppendFloat(t.output, float64(*value), format, precision, 32)
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayFloat64(value *float64, format byte, precision int) error {
	t.appendDelimiter()
	if value != nil {
		t.output = strconv.AppendFloat(t.output, *value, format, precision, 64)
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayString(value *string) error {
	t.appendDelimiter()
	if value != nil {
		t.append([]byte(*value))
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayByteString(value *[]byte) error {
	t.appendDelimiter()
	if value != nil {
		t.append(*value)
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayBytes(value *[]byte) error {
	t.appendDelimiter()
	if value != nil {
		outLen := len(t.output)
		t.output = simpleutils.ResizeBytes(t.output, outLen+base64.StdEncoding.EncodedLen(len(*value)))
		base64.StdEncoding.Encode(t.output[outLen:], *value)
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayBool(value *bool) error {
	t.appendDelimiter()
	if value != nil {
		t.output = strconv.AppendBool(t.output, *value)
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayBigInt(value *big.Int) error {
	t.appendDelimiter()
	if value != nil {
		data, _ := value.MarshalText()
		t.append(data)
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayRat(value *big.Rat, precision int) error {
	t.appendDelimiter()
	if value != nil {
		t.append([]byte(value.FloatString(precision)))
	}
	return nil
}
func (t *ToDelimitedBytesImpl) ArrayBigFloat(value *big.Float, format byte, precision int) error {
	t.appendDelimiter()
	if value != nil {
		t.append([]byte(value.Text(format, precision)))
	}
	return nil
}

func (t *ToDelimitedBytesImpl) ValueInt(value *int, typeName string, typeDescription string) error {
	return t.ArrayInt(value)
}
func (t *ToDelimitedBytesImpl) ValueInt32(value *int32, typeName string, typeDescription string) error {
	return t.ArrayInt32(value)
}
func (t *ToDelimitedBytesImpl) ValueInt64(value *int64, typeName string, typeDescription string) error {
	return t.ArrayInt64(value)
}
func (t *ToDelimitedBytesImpl) ValueFloat32(value *float32, format byte, precision int, typeName string, typeDescription string) error {
	return t.ArrayFloat32(value, format, precision)
}
func (t *ToDelimitedBytesImpl) ValueFloat64(value *float64, format byte, precision int, typeName string, typeDescription string) error {
	return t.ArrayFloat64(value, format, precision)
}
func (t *ToDelimitedBytesImpl) ValueString(value *string, typeName string, typeDescription string) error {
	return t.ArrayString(value)
}
func (t *ToDelimitedBytesImpl) ValueByteString(value *[]byte, typeName string, typeDescription string) error {
	return t.ArrayByteString(value)
}
func (t *ToDelimitedBytesImpl) ValueBytes(value *[]byte, typeName string, typeDescription string) error {
	return t.ArrayBytes(value)
}
func (t *ToDelimitedBytesImpl) ValueBool(value *bool, typeName string, typeDescription string) error {
	return t.ArrayBool(value)
}
func (t *ToDelimitedBytesImpl) ValueBigInt(value *big.Int, typeName string, typeDescription string) error {
	return t.ArrayBigInt(value)
}
func (t *ToDelimitedBytesImpl) ValueRat(value *big.Rat, precision int, typeName string, typeDescription string) error {
	return t.ArrayRat(value, precision)
}
func (t *ToDelimitedBytesImpl) ValueBigFloat(value *big.Float, format byte, precision int, typeName string, typeDescription string) error {
	return t.ArrayBigFloat(value, format, precision)
}

func (t *ToDelimitedBytesImpl) StartObject(typeName string, typeDescription string) error {
	return nil
}
func (t *ToDelimitedBytesImpl) StartArray(typeName string, valueTypeName string, typeDescription string) error {
	return nil
}
func (t *ToDelimitedBytesImpl) StartMap(typeName string, valueTypeName string, typeDescription string) error {
	return nil
}
func (t *ToDelimitedBytesImpl) EndObject() error {
	return nil
}
func (t *ToDelimitedBytesImpl) EndArray() error {
	return nil
}
func (t *ToDelimitedBytesImpl) EndMap() error {
	return nil
}

func (t *ToDelimitedBytesImpl) ObjectInt(value *int, name string, mandatory bool, description string) error {
	return t.ArrayInt(value)
}
func (t *ToDelimitedBytesImpl) ObjectInt32(value *int32, name string, mandatory bool, description string) error {
	return t.ArrayInt32(value)
}
func (t *ToDelimitedBytesImpl) ObjectInt64(value *int64, name string, mandatory bool, description string) error {
	return t.ArrayInt64(value)
}
func (t *ToDelimitedBytesImpl) ObjectFloat32(value *float32, format byte, precision int, name string, mandatory bool, description string) error {
	return t.ArrayFloat32(value, format, precision)
}
func (t *ToDelimitedBytesImpl) ObjectFloat64(value *float64, format byte, precision int, name string, mandatory bool, description string) error {
	return t.ArrayFloat64(value, format, precision)
}
func (t *ToDelimitedBytesImpl) ObjectString(value *string, name string, mandatory bool, description string) error {
	return t.ArrayString(value)
}
func (t *ToDelimitedBytesImpl) ObjectByteString(value *[]byte, name string, mandatory bool, description string) error {
	return t.ArrayByteString(value)
}
func (t *ToDelimitedBytesImpl) ObjectBytes(value *[]byte, name string, mandatory bool, description string) error {
	return t.ArrayBytes(value)
}
func (t *ToDelimitedBytesImpl) ObjectBool(value *bool, name string, mandatory bool, description string) error {
	return t.ArrayBool(value)
}
func (t *ToDelimitedBytesImpl) ObjectBigInt(value *big.Int, name string, mandatory bool, description string) error {
	return t.ArrayBigInt(value)
}
func (t *ToDelimitedBytesImpl) ObjectRat(value *big.Rat, precision int, name string, mandatory bool, description string) error {
	return t.ArrayRat(value, precision)
}
func (t *ToDelimitedBytesImpl) ObjectBigFloat(value *big.Float, format byte, precision int, name string, mandatory bool, description string) error {
	return t.ArrayBigFloat(value, format, precision)
}
func (t *ToDelimitedBytesImpl) ObjectValue(name string, mandatory bool, description string) (bool, error) {
	return true, nil
}

func (t *ToDelimitedBytesImpl) ArrayLen(length int) (int, error) {
	return 0, t.ArrayInt(&length)
}
func (t *ToDelimitedBytesImpl) ArrayValue() error {
	return nil
}

func (t *ToDelimitedBytesImpl) MapLen(length int) (int, error) {
	return 0, t.ArrayInt(&length)
}
func (t *ToDelimitedBytesImpl) MapNextKey() (string, error) {
	return "", nil
}
func (t *ToDelimitedBytesImpl) MapReadInt() (int, error) {
	return 0, nil
}
func (t *ToDelimitedBytesImpl) MapReadInt32() (int32, error) {
	return 0, nil
}
func (t *ToDelimitedBytesImpl) MapReadInt64() (int64, error) {
	return 0, nil
}
func (t *ToDelimitedBytesImpl) MapReadFloat32() (float32, error) {
	return 0, nil
}
func (t *ToDelimitedBytesImpl) MapReadFloat64() (float64, error) {
	return 0, nil
}
func (t *ToDelimitedBytesImpl) MapReadString() (string, error) {
	return "", nil
}
func (t *ToDelimitedBytesImpl) MapReadByteString(value []byte) ([]byte, error) {
	return nil, nil
}
func (t *ToDelimitedBytesImpl) MapReadBytes(value []byte) ([]byte, error) {
	return nil, nil
}
func (t *ToDelimitedBytesImpl) MapReadBool() (bool, error) {
	return false, nil
}
func (t *ToDelimitedBytesImpl) MapReadBigInt(value *big.Int) (*big.Int, error) {
	return value, nil
}
func (t *ToDelimitedBytesImpl) MapReadRat(value *big.Rat) (*big.Rat, error) {
	return value, nil
}
func (t *ToDelimitedBytesImpl) MapReadBigFloat(value *big.Float) (*big.Float, error) {
	return value, nil
}
func (t *ToDelimitedBytesImpl) MapReadValue() error {
	return nil
}

func (t *ToDelimitedBytesImpl) AppendKey(key string) {
	t.appendDelimiter()
	t.append([]byte(key))
}

func (t *ToDelimitedBytesImpl) MapWriteInt(key string, value int) error {
	t.AppendKey(key)
	return t.ArrayInt(&value)
}
func (t *ToDelimitedBytesImpl) MapWriteInt32(key string, value int32) error {
	t.AppendKey(key)
	return t.ArrayInt32(&value)
}
func (t *ToDelimitedBytesImpl) MapWriteInt64(key string, value int64) error {
	t.AppendKey(key)
	return t.ArrayInt64(&value)
}
func (t *ToDelimitedBytesImpl) MapWriteFloat32(key string, value float32, format byte, precision int) error {
	t.AppendKey(key)
	return t.ArrayFloat32(&value, format, precision)
}
func (t *ToDelimitedBytesImpl) MapWriteFloat64(key string, value float64, format byte, precision int) error {
	t.AppendKey(key)
	return t.ArrayFloat64(&value, format, precision)
}
func (t *ToDelimitedBytesImpl) MapWriteString(key string, value string) error {
	t.AppendKey(key)
	return t.ArrayString(&value)
}
func (t *ToDelimitedBytesImpl) MapWriteByteString(key string, value []byte) error {
	t.AppendKey(key)
	return t.ArrayByteString(&value)
}
func (t *ToDelimitedBytesImpl) MapWriteBytes(key string, value []byte) error {
	t.AppendKey(key)
	return t.ArrayBytes(&value)
}
func (t *ToDelimitedBytesImpl) MapWriteBool(key string, value bool) error {
	t.AppendKey(key)
	return t.ArrayBool(&value)
}
func (t *ToDelimitedBytesImpl) MapWriteBigInt(key string, value *big.Int) error {
	t.AppendKey(key)
	return t.ArrayBigInt(value)
}
func (t *ToDelimitedBytesImpl) MapWriteRat(key string, value *big.Rat, precision int) error {
	t.AppendKey(key)
	return t.ArrayRat(value, precision)
}
func (t *ToDelimitedBytesImpl) MapWriteBigFloat(key string, value *big.Float, format byte, precision int) error {
	t.AppendKey(key)
	return t.ArrayBigFloat(value, format, precision)
}
func (t *ToDelimitedBytesImpl) MapWriteValue(key string) error {
	t.AppendKey(key)
	return nil
}

func (t *ToDelimitedBytesImpl) Clear() {
	t.output = t.output[:0]
	t.hadAnything = false
}

func (t *ToDelimitedBytesImpl) Results() []byte {
	return t.output
}

func init() {
	var _ inspect.InspectorImpl = (*ToDelimitedBytesImpl)(nil)
}
