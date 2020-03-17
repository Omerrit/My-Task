package fromjson

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/buger/jsonparser"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectorhelpers"
	"gerrit-share.lan/go/utils/simpleutils"
	"math/big"
	"strconv"
	"unsafe"
)

var trueValue = []byte("true")
var falseValue = []byte("false")

type value struct {
	buffer    []byte
	offset    int
	valueType jsonparser.ValueType
}

type object map[string]value

type array struct {
	values []value
	index  int
}

type mapelement struct {
	key   string
	value value
}

type mapArray struct {
	elems []mapelement
	index int
}

func (v *value) makeError(err error) error {
	//panic(&fullParseError{err, v.offset})
	return &fullParseError{err, v.offset}
}

func (v *value) makePropertyNotFoundError(name string) error {
	return &fullParseError{fmt.Errorf("%w: %s", ErrPropertyNotFound, name), v.offset}
}

func (v *value) toString() string {
	return *(*string)(unsafe.Pointer(&v.buffer))
}

func (v *value) checkNull() error {
	if !isNull(v.buffer) {
		return &fullParseError{ErrShouldBeNull, v.offset}
	}
	return nil
}

type Inspector struct {
	inspectorhelpers.Reader
	values          []value
	currentValue    value
	objects         []object
	currentObject   object
	arrays          []array
	currentArray    array
	mapArrays       []mapArray
	currentMapArray mapArray
}

func (i *Inspector) popValue() {
	if len(i.values) == 0 {
		return
	}
	i.currentValue = i.values[len(i.values)-1]
	i.values = i.values[:(len(i.values) - 1)]
}

func (i *Inspector) ValueInt(value *int, typeName string, typeDescription string) error {
	var err error
	if value == nil {
		err = i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.Number && i.currentValue.valueType != jsonparser.Unknown {
		err = i.currentValue.makeError(ErrIntegerRequired)
		return err
	}
	*value, err = strconv.Atoi(i.currentValue.toString())
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	i.popValue()
	return nil
}

func (i *Inspector) ValueInt32(value *int32, typeName string, typeDescription string) error {
	if value == nil {
		err := i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.Number && i.currentValue.valueType != jsonparser.Unknown {
		err := i.currentValue.makeError(ErrIntegerRequired)
		return err
	}
	v, err := strconv.ParseInt(i.currentValue.toString(), 10, 32)
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	*value = int32(v)
	i.popValue()
	return nil
}

func (i *Inspector) ValueInt64(value *int64, typeName string, typeDescription string) error {
	var err error
	if value == nil {
		err = i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.Number && i.currentValue.valueType != jsonparser.Unknown {
		err = i.currentValue.makeError(ErrIntegerRequired)
		return err
	}
	*value, err = strconv.ParseInt(i.currentValue.toString(), 10, 64)
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	i.popValue()
	return nil
}

func (i *Inspector) ValueFloat32(value *float32, format byte, precision int, typeName string, typeDescription string) error {
	if value == nil {
		err := i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.Number && i.currentValue.valueType != jsonparser.Unknown {
		err := i.currentValue.makeError(ErrFloatRequired)
		return err
	}
	v, err := strconv.ParseFloat(i.currentValue.toString(), 32)
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	*value = float32(v)
	i.popValue()
	return nil
}

func (i *Inspector) ValueFloat64(value *float64, format byte, precision int, typeName string, typeDescription string) error {
	var err error
	if value == nil {
		err = i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.Number && i.currentValue.valueType != jsonparser.Unknown {
		err = i.currentValue.makeError(ErrFloatRequired)
		return err
	}
	*value, err = strconv.ParseFloat(i.currentValue.toString(), 64)
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	i.popValue()
	return nil
}

func (i *Inspector) ValueString(value *string, typeName string, typeDescription string) error {
	var err error
	if value == nil {
		err = i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.String && i.currentValue.valueType != jsonparser.Unknown {
		err = i.currentValue.makeError(ErrStringRequired)
		return err
	}
	var unescaped []byte
	unescaped, err = jsonparser.Unescape(i.currentValue.buffer, nil)
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	*value = string(unescaped)
	i.popValue()
	return nil
}

func (i *Inspector) ValueByteString(value *[]byte, typeName string, typeDescription string) error {
	var err error
	if value == nil {
		err = i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.String && i.currentValue.valueType != jsonparser.Unknown {
		err = i.currentValue.makeError(ErrStringRequired)
		return err
	}
	*value, err = jsonparser.Unescape(i.currentValue.buffer, *value)
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	i.popValue()
	return nil
}

func (i *Inspector) ValueBytes(value *[]byte, typeName string, typeDescription string) error {
	if value == nil {
		err := i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.String && i.currentValue.valueType != jsonparser.Unknown {
		err := i.currentValue.makeError(ErrStringRequired)
		return err
	}
	*value = simpleutils.ResizeBytes(*value, base64.RawURLEncoding.DecodedLen(len(i.currentValue.buffer)))
	n, err := base64.RawURLEncoding.Decode(*value, i.currentValue.buffer)
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	*value = (*value)[:n]
	i.popValue()
	return nil
}

func (i *Inspector) ValueBool(value *bool, typeName string, typeDescription string) error {
	if value == nil {
		err := i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.Boolean && i.currentValue.valueType != jsonparser.Unknown {
		err := i.currentValue.makeError(ErrBoolRequired)
		return err
	}
	if bytes.Equal(i.currentValue.buffer, trueValue) {
		*value = true
	} else if bytes.Equal(i.currentValue.buffer, falseValue) {
		*value = false
	} else {
		err := i.currentValue.makeError(ErrNotBool)
		return err
	}
	i.popValue()
	return nil
}

func (i *Inspector) ValueBigInt(value *big.Int, typeName string, typeDescription string) error {
	if value == nil {
		err := i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.String && i.currentValue.valueType != jsonparser.Unknown {
		err := i.currentValue.makeError(ErrBigIntRequired)
		return err
	}
	err := value.UnmarshalText(i.currentValue.buffer)
	if err != nil {
		err := i.currentValue.makeError(err)
		return err
	}
	i.popValue()
	return nil
}

func (i *Inspector) ValueRat(value *big.Rat, precision int, typeName string, typeDescription string) error {
	if value == nil {
		err := i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.String && i.currentValue.valueType != jsonparser.Unknown {
		err := i.currentValue.makeError(ErrRatRequired)
		return err
	}
	_, ok := value.SetString(i.currentValue.toString())
	if !ok {
		err := i.currentValue.makeError(ErrParsingValue)
		return err
	}
	i.popValue()
	return nil
}

func (i *Inspector) ValueBigFloat(value *big.Float, format byte, precision int, typeName string, typeDescription string) error {
	if value == nil {
		err := i.currentValue.checkNull()
		return err
	}
	if i.currentValue.valueType != jsonparser.String && i.currentValue.valueType != jsonparser.Unknown {
		err := i.currentValue.makeError(ErrBigFloatRequired)
		return err
	}
	_, ok := value.SetString(i.currentValue.toString())
	if !ok {
		err := i.currentValue.makeError(ErrParsingValue)
		return err
	}
	i.popValue()
	return nil
}

func (i *Inspector) StartObject(typeName string, typeDescription string) error {
	if i.currentValue.valueType != jsonparser.Object && i.currentValue.valueType != jsonparser.Unknown {
		err := i.currentValue.makeError(ErrObjectRequired)
		return err
	}
	object := make(map[string]value)
	err := jsonparser.ObjectEach(i.currentValue.buffer, func(key []byte, data []byte, valueType jsonparser.ValueType, offset int) error {
		object[string(key)] = value{data, offset + i.currentValue.offset, valueType}
		return nil
	})
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	i.objects = append(i.objects, i.currentObject)
	i.currentObject = object
	return nil
}

func (i *Inspector) StartArray(typeName string, valueTypeName string, typeDescription string) error {
	if i.currentValue.valueType != jsonparser.Array && i.currentValue.valueType != jsonparser.Unknown {
		err := i.currentValue.makeError(ErrArrayRequired)
		return err
	}
	array := array{}
	var parseError error
	_, err := jsonparser.ArrayEach(i.currentValue.buffer, func(data []byte, valueType jsonparser.ValueType, offset int, err error) {
		array.values = append(array.values, value{data, offset + i.currentValue.offset, valueType})
		if parseError == nil && err != nil {
			parseError = makeParseError(offset, err)
		}
	})
	if parseError != nil {
		return parseError
	}
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	i.arrays = append(i.arrays, i.currentArray)
	i.currentArray = array
	return nil
}

func (i *Inspector) StartMap(typeName string, valueTypeName string, typeDescription string) error {
	if i.currentValue.valueType != jsonparser.Object && i.currentValue.valueType != jsonparser.Unknown {
		err := i.currentValue.makeError(ErrObjectRequired)
		return err
	}
	mapArray := mapArray{}
	err := jsonparser.ObjectEach(i.currentValue.buffer, func(key []byte, data []byte, valueType jsonparser.ValueType, offset int) error {
		mapArray.elems = append(mapArray.elems, mapelement{string(key), value{data, offset + i.currentValue.offset, valueType}})
		return nil
	})
	if err != nil {
		err = i.currentValue.makeError(err)
		return err
	}
	i.mapArrays = append(i.mapArrays, i.currentMapArray)
	i.currentMapArray = mapArray
	return nil
}

func (i *Inspector) EndObject() error {
	if len(i.objects) == 0 {
		return nil
	}
	i.currentObject = i.objects[len(i.objects)-1]
	i.objects = i.objects[:(len(i.objects) - 1)]
	i.popValue()
	return nil
}

func (i *Inspector) EndArray() error {
	if len(i.arrays) == 0 {
		return nil
	}
	i.currentArray = i.arrays[len(i.arrays)-1]
	i.arrays = i.arrays[:(len(i.arrays) - 1)]
	i.popValue()
	return nil
}

func (i *Inspector) EndMap() error {
	if len(i.mapArrays) == 0 {
		return nil
	}
	i.currentMapArray = i.mapArrays[len(i.mapArrays)-1]
	i.mapArrays = i.mapArrays[:(len(i.mapArrays) - 1)]
	i.popValue()
	return nil
}

func (i *Inspector) ObjectInt(value *int, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Number {
		return v.makeError(ErrIntegerRequired)
	}
	var err error
	*value, err = strconv.Atoi(v.toString())
	if err != nil {
		return v.makeError(err)
	}
	return nil
}

func (i *Inspector) ObjectInt32(value *int32, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Number {
		return v.makeError(ErrIntegerRequired)
	}
	result, err := strconv.ParseInt(v.toString(), 10, 32)
	if err != nil {
		return v.makeError(err)
	}
	*value = int32(result)
	return nil
}

func (i *Inspector) ObjectInt64(value *int64, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Number {
		return v.makeError(ErrIntegerRequired)
	}
	var err error
	*value, err = strconv.ParseInt(v.toString(), 10, 64)
	if err != nil {
		return v.makeError(err)
	}
	return nil
}

func (i *Inspector) ObjectFloat32(value *float32, format byte, precision int, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Number {
		return v.makeError(ErrIntegerRequired)
	}
	result, err := strconv.ParseFloat(v.toString(), 32)
	if err != nil {
		return v.makeError(err)
	}
	*value = float32(result)
	return nil
}

func (i *Inspector) ObjectFloat64(value *float64, format byte, precision int, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Number {
		return v.makeError(ErrIntegerRequired)
	}
	var err error
	*value, err = strconv.ParseFloat(v.toString(), 64)
	if err != nil {
		return v.makeError(err)
	}
	return nil
}

func (i *Inspector) ObjectString(value *string, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrStringRequired)
	}
	unescaped, err := jsonparser.Unescape(v.buffer, nil)
	if err != nil {
		return v.makeError(err)
	}
	*value = string(unescaped)
	return nil
}

func (i *Inspector) ObjectByteString(value *[]byte, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrStringRequired)
	}
	var err error
	*value, err = jsonparser.Unescape(v.buffer, *value)
	if err != nil {
		return v.makeError(err)
	}
	return nil
}

func (i *Inspector) ObjectBytes(value *[]byte, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrStringRequired)
	}
	*value = simpleutils.ResizeBytes(*value, base64.RawURLEncoding.DecodedLen(len(v.buffer)))
	n, err := base64.RawURLEncoding.Decode(*value, v.buffer)
	if err != nil {
		return v.makeError(err)
	}
	*value = (*value)[:n]
	return nil
}

func (i *Inspector) ObjectBool(value *bool, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Boolean {
		return v.makeError(ErrStringRequired)
	}
	if bytes.Equal(v.buffer, trueValue) {
		*value = true
	} else if bytes.Equal(v.buffer, falseValue) {
		*value = false
	} else {
		return v.makeError(ErrNotBool)
	}
	return nil
}

func (i *Inspector) ObjectBigInt(value *big.Int, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrBigIntRequired)
	}
	err := value.UnmarshalText(v.buffer)
	if err != nil {
		return v.makeError(err)
	}
	return nil
}

func (i *Inspector) ObjectRat(value *big.Rat, precision int, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrRatRequired)
	}
	_, ok = value.SetString(v.toString())
	if !ok {
		return v.makeError(ErrParsingValue)
	}
	return nil
}

func (i *Inspector) ObjectBigFloat(value *big.Float, format byte, precision int, name string, mandatory bool, description string) error {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return i.currentValue.makePropertyNotFoundError(name)
		}
		return nil
	}
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrBigFloatRequired)
	}
	_, ok = value.SetString(v.toString())
	if !ok {
		return v.makeError(ErrParsingValue)
	}
	return nil
}

func (i *Inspector) ObjectValue(name string, mandatory bool, description string) (bool, error) {
	v, ok := i.currentObject[name]
	if !ok {
		if mandatory {
			return true, i.currentValue.makePropertyNotFoundError(name)
		}
		return false, nil
	}
	i.values = append(i.values, i.currentValue)
	i.currentValue = v
	return true, nil
}

func (i *Inspector) ArrayLen(length int) (int, error) {
	return len(i.currentArray.values), nil
}

func (i *Inspector) ArrayInt(value *int) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Number {
		return v.makeError(ErrIntegerRequired)
	}
	var err error
	*value, err = strconv.Atoi(v.toString())
	if err != nil {
		return v.makeError(err)
	}
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayInt32(value *int32) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Number {
		return v.makeError(ErrIntegerRequired)
	}
	result, err := strconv.ParseInt(v.toString(), 10, 32)
	if err != nil {
		return v.makeError(err)
	}
	*value = int32(result)
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayInt64(value *int64) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Number {
		return v.makeError(ErrIntegerRequired)
	}
	var err error
	*value, err = strconv.ParseInt(v.toString(), 10, 64)
	if err != nil {
		return v.makeError(err)
	}
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayFloat32(value *float32, format byte, precision int) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Number {
		return v.makeError(ErrFloatRequired)
	}
	result, err := strconv.ParseFloat(v.toString(), 32)
	if err != nil {
		return v.makeError(err)
	}
	*value = float32(result)
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayFloat64(value *float64, format byte, precision int) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.Number {
		return v.makeError(ErrFloatRequired)
	}
	var err error
	*value, err = strconv.ParseFloat(v.toString(), 64)
	if err != nil {
		return v.makeError(err)
	}
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayString(value *string) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrStringRequired)
	}
	unescaped, err := jsonparser.Unescape(v.buffer, nil)
	if err != nil {
		return v.makeError(err)
	}
	*value = string(unescaped)
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayByteString(value *[]byte) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrStringRequired)
	}
	var err error
	*value, err = jsonparser.Unescape(v.buffer, *value)
	if err != nil {
		return v.makeError(err)
	}
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayBytes(value *[]byte) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrStringRequired)
	}
	*value = simpleutils.ResizeBytes(*value, base64.RawURLEncoding.DecodedLen(len(v.buffer)))
	n, err := base64.RawURLEncoding.Decode(*value, v.buffer)
	if err != nil {
		return v.makeError(err)
	}
	*value = (*value)[:n]
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayBool(value *bool) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrBoolRequired)
	}
	if bytes.Equal(v.buffer, trueValue) {
		*value = true
	} else if bytes.Equal(v.buffer, falseValue) {
		*value = false
	} else {
		return v.makeError(ErrNotBool)
	}
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayBigInt(value *big.Int) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrBigIntRequired)
	}
	err := value.UnmarshalText(v.buffer)
	if err != nil {
		return v.makeError(err)
	}
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayRat(value *big.Rat, precision int) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrRatRequired)
	}
	_, ok := value.SetString(v.toString())
	if !ok {
		return v.makeError(ErrParsingValue)
	}
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayBigFloat(value *big.Float, format byte, precision int) error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	v := &i.currentArray.values[i.currentArray.index]
	if value == nil {
		return v.checkNull()
	}
	if v.valueType != jsonparser.String {
		return v.makeError(ErrBigFloatRequired)
	}
	_, ok := value.SetString(v.toString())
	if !ok {
		return v.makeError(ErrParsingValue)
	}
	i.currentArray.index++
	return nil
}

func (i *Inspector) ArrayValue() error {
	if i.currentArray.index >= len(i.currentArray.values) {
		return i.currentValue.makeError(ErrWrongLength)
	}
	i.values = append(i.values, i.currentValue)
	i.currentValue = i.currentArray.values[i.currentArray.index]
	i.currentArray.index++
	return nil
}

func (i *Inspector) MapLen(length int) (int, error) {
	return len(i.currentMapArray.elems), nil
}

func (i *Inspector) MapNextKey() (string, error) {
	if i.currentMapArray.index >= len(i.currentMapArray.elems) {
		return "", i.currentValue.makeError(ErrWrongLength)
	}
	key := i.currentMapArray.elems[i.currentMapArray.index].key
	return key, nil
}

func (i *Inspector) MapReadInt() (int, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.Number {
		return 0, v.makeError(ErrIntegerRequired)
	}
	result, err := strconv.Atoi(v.toString())
	if err != nil {
		return 0, v.makeError(err)
	}
	i.currentMapArray.index++
	return result, nil
}

func (i *Inspector) MapReadInt32() (int32, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.Number {
		return 0, v.makeError(ErrIntegerRequired)
	}
	result, err := strconv.ParseInt(v.toString(), 10, 32)
	if err != nil {
		return 0, v.makeError(err)
	}
	i.currentMapArray.index++
	return int32(result), nil
}

func (i *Inspector) MapReadInt64() (int64, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.Number {
		return 0, v.makeError(ErrIntegerRequired)
	}
	result, err := strconv.ParseInt(v.toString(), 10, 64)
	if err != nil {
		return 0, v.makeError(err)
	}
	i.currentMapArray.index++
	return result, nil
}

func (i *Inspector) MapReadFloat32() (float32, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.Number {
		return 0, v.makeError(ErrFloatRequired)
	}
	result, err := strconv.ParseFloat(v.toString(), 32)
	if err != nil {
		return 0, v.makeError(err)
	}
	i.currentMapArray.index++
	return float32(result), nil
}

func (i *Inspector) MapReadFloat64() (float64, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.Number {
		return 0, v.makeError(ErrFloatRequired)
	}
	result, err := strconv.ParseFloat(v.toString(), 64)
	if err != nil {
		return 0, v.makeError(err)
	}
	i.currentMapArray.index++
	return result, nil
}

func (i *Inspector) MapReadString() (string, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.String {
		return "", v.makeError(ErrStringRequired)
	}
	result, err := jsonparser.Unescape(v.buffer, nil)
	if err != nil {
		return "", v.makeError(err)
	}
	i.currentMapArray.index++
	return string(result), nil
}

func (i *Inspector) MapReadByteString(value []byte) ([]byte, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.String {
		return value, v.makeError(ErrStringRequired)
	}
	var err error
	value, err = jsonparser.Unescape(v.buffer, value)
	if err != nil {
		return value, v.makeError(err)
	}
	i.currentMapArray.index++
	return value, nil
}

func (i *Inspector) MapReadBytes(value []byte) ([]byte, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.String {
		return value, v.makeError(ErrStringRequired)
	}
	value = simpleutils.ResizeBytes(value, base64.RawURLEncoding.DecodedLen(len(v.buffer)))
	n, err := base64.RawURLEncoding.Decode(value, v.buffer)
	if err != nil {
		return value, v.makeError(err)
	}
	i.currentMapArray.index++
	value = value[:n]
	return value, nil
}

func (i *Inspector) MapReadBool() (bool, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.Boolean {
		return false, v.makeError(ErrBoolRequired)
	}
	if bytes.Equal(v.buffer, trueValue) {
		i.currentMapArray.index++
		return true, nil
	} else if bytes.Equal(v.buffer, falseValue) {
		i.currentMapArray.index++
		return false, nil
	} else {
		return false, v.makeError(ErrNotBool)
	}
}

func (i *Inspector) MapReadBigInt(value *big.Int) (*big.Int, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.String {
		return value, v.makeError(ErrBigIntRequired)
	}
	err := value.UnmarshalText(v.buffer)
	if err != nil {
		return nil, v.makeError(err)
	}
	i.currentMapArray.index++
	return value, nil
}

func (i *Inspector) MapReadRat(value *big.Rat) (*big.Rat, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.String {
		return value, v.makeError(ErrRatRequired)
	}
	_, ok := value.SetString(v.toString())
	if !ok {
		return value, v.makeError(ErrParsingValue)
	}
	i.currentMapArray.index++
	return value, nil
}

func (i *Inspector) MapReadBigFloat(value *big.Float) (*big.Float, error) {
	v := &i.currentMapArray.elems[i.currentMapArray.index].value
	if v.valueType != jsonparser.String {
		return value, v.makeError(ErrBigFloatRequired)
	}
	_, ok := value.SetString(v.toString())
	if !ok {
		return value, v.makeError(ErrParsingValue)
	}
	i.currentMapArray.index++
	return value, nil
}

func (i *Inspector) MapReadValue() error {
	i.values = append(i.values, i.currentValue)
	i.currentValue = i.currentMapArray.elems[i.currentMapArray.index].value
	i.currentMapArray.index++
	return nil
}

func NewInspector(data []byte, dataOffset int) *Inspector {
	return &Inspector{currentValue: value{data, dataOffset, jsonparser.Unknown}}
}

func init() {
	var _ inspect.InspectorImpl = (*Inspector)(nil)
}
