package serializers

import (
	"encoding/base64"
	"fmt"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectorhelpers"
	"gerrit-share.lan/go/utils/multisets"
	"math/big"
	"net/url"
	"strconv"
)

type FromUrl struct {
	inspectorhelpers.Reader
	Values     url.Values
	counts     multisets.String
	currentStr string
	inProgress bool
}

func (i *FromUrl) getString(name string, mandatory bool) (str string, skip bool, err error) {
	strs, ok := i.Values[name]
	if !ok {
		if mandatory {
			return "", true, fmt.Errorf("%w: %s", inspect.ErrMandatoryFieldAbsent, name)
		} else {
			return "", true, nil
		}
	}
	offset := i.counts.Count(name)
	if len(strs) <= offset {
		if mandatory {
			return "", true, fmt.Errorf("%w: %s #%d, have: %v", inspect.ErrMandatoryFieldAbsent, name, offset, strs)
		} else {
			return "", true, nil
		}
	}
	i.counts.Add(name)
	return strs[offset], false, nil
}

func (i *FromUrl) ValueInt(value *int, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	var err error
	*value, err = strconv.Atoi(i.currentStr)
	i.currentStr = ""
	return err
}

func (i *FromUrl) ValueInt32(value *int32, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	var err error
	var v int64
	v, err = strconv.ParseInt(i.currentStr, 10, 32)
	*value = int32(v)
	i.currentStr = ""
	return err
}

func (i *FromUrl) ValueInt64(value *int64, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	var err error
	*value, err = strconv.ParseInt(i.currentStr, 10, 64)
	i.currentStr = ""
	return err
}

func (i *FromUrl) ValueFloat32(value *float32, format byte, precision int, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	var err error
	var v float64
	v, err = strconv.ParseFloat(i.currentStr, 32)
	*value = float32(v)
	i.currentStr = ""
	return err
}

func (i *FromUrl) ValueFloat64(value *float64, format byte, precision int, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	var err error
	*value, err = strconv.ParseFloat(i.currentStr, 64)
	i.currentStr = ""
	return err
}

func (i *FromUrl) ValueString(value *string, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	*value = i.currentStr
	i.currentStr = ""
	return nil
}

func (i *FromUrl) ValueByteString(value *[]byte, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	*value = []byte(i.currentStr)
	i.currentStr = ""
	return nil
}

func (i *FromUrl) ValueBytes(value *[]byte, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	var err error
	*value, err = base64.StdEncoding.DecodeString(i.currentStr)
	i.currentStr = ""
	return err
}

func (i *FromUrl) ValueBool(value *bool, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	var err error
	*value, err = strconv.ParseBool(i.currentStr)
	i.currentStr = ""
	return err
}

func (i *FromUrl) ValueBigInt(value *big.Int, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	err := value.UnmarshalText([]byte(i.currentStr))
	i.currentStr = ""
	return err
}

func (i *FromUrl) ValueRat(value *big.Rat, precision int, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	err := value.UnmarshalText([]byte(i.currentStr))
	i.currentStr = ""
	return err
}

func (i *FromUrl) ValueBigFloat(value *big.Float, format byte, precision int, typeName string, typeDescription string) error {
	if len(i.currentStr) == 0 {
		return nil
	}
	err := value.UnmarshalText([]byte(i.currentStr))
	i.currentStr = ""
	return err
}

func (i *FromUrl) StartObject(typeName string, typeDescription string) error {
	if i.inProgress {
		return ErrUnsupportedType
	}
	i.inProgress = true
	return nil
}

func (i *FromUrl) StartArray(typeName string, valueTypeName string, typeDescription string) error {
	return ErrUnsupportedType
}

func (i *FromUrl) StartMap(typeName string, valueTypeName string, typeDescription string) error {
	return ErrUnsupportedType
}

func (i *FromUrl) EndObject() error {

	return nil
}

func (i *FromUrl) EndArray() error {

	return nil
}

func (i *FromUrl) EndMap() error {
	return nil
}

func (i *FromUrl) ObjectInt(value *int, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	*value, err = strconv.Atoi(str)
	return err
}

func (i *FromUrl) ObjectInt32(value *int32, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	var v int64
	v, err = strconv.ParseInt(str, 10, 32)
	*value = int32(v)
	return err
}

func (i *FromUrl) ObjectInt64(value *int64, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	*value, err = strconv.ParseInt(str, 10, 64)
	return nil
}

func (i *FromUrl) ObjectFloat32(value *float32, format byte, precision int, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	var v float64
	v, err = strconv.ParseFloat(str, 32)
	*value = float32(v)
	return err
}

func (i *FromUrl) ObjectFloat64(value *float64, format byte, precision int, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	*value, err = strconv.ParseFloat(str, 64)
	return err
}

func (i *FromUrl) ObjectString(value *string, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	*value = str
	return nil
}

func (i *FromUrl) ObjectByteString(value *[]byte, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	*value = []byte(str)
	return nil
}

func (i *FromUrl) ObjectBytes(value *[]byte, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	*value, err = base64.StdEncoding.DecodeString(str)
	return err
}

func (i *FromUrl) ObjectBool(value *bool, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	*value, err = strconv.ParseBool(str)
	return err
}

func (i *FromUrl) ObjectBigInt(value *big.Int, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	return value.UnmarshalText([]byte(str))
}

func (i *FromUrl) ObjectRat(value *big.Rat, precision int, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	return value.UnmarshalText([]byte(str))
}

func (i *FromUrl) ObjectBigFloat(value *big.Float, format byte, precision int, name string, mandatory bool, description string) error {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return err
	}
	return value.UnmarshalText([]byte(str))
}

func (i *FromUrl) ObjectValue(name string, mandatory bool, description string) (bool, error) {
	str, skip, err := i.getString(name, mandatory)
	if skip {
		return !skip, err
	}
	i.currentStr = str
	return true, nil
}

func (i *FromUrl) ArrayLen(length int) (int, error) {
	return length, ErrUnsupportedType
}

func (i *FromUrl) ArrayInt(value *int) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayInt32(value *int32) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayInt64(value *int64) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayFloat32(value *float32, format byte, precision int) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayFloat64(value *float64, format byte, precision int) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayString(value *string) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayByteString(value *[]byte) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayBytes(value *[]byte) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayBool(value *bool) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayBigInt(value *big.Int) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayRat(value *big.Rat, precision int) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayBigFloat(value *big.Float, format byte, precision int) error {
	return ErrUnsupportedType
}

func (i *FromUrl) ArrayValue() error {
	return ErrUnsupportedType
}

func (i *FromUrl) MapLen(length int) (int, error) {
	return length, ErrUnsupportedType
}

func (i *FromUrl) MapNextKey() (string, error) {
	return "", ErrUnsupportedType
}

func (i *FromUrl) MapReadInt() (int, error) {
	return 0, ErrUnsupportedType
}

func (i *FromUrl) MapReadInt32() (int32, error) {
	return 0, ErrUnsupportedType
}

func (i *FromUrl) MapReadInt64() (int64, error) {
	return 0, ErrUnsupportedType
}

func (i *FromUrl) MapReadFloat32() (float32, error) {
	return 0, ErrUnsupportedType
}

func (i *FromUrl) MapReadFloat64() (float64, error) {
	return 0, ErrUnsupportedType
}

func (i *FromUrl) MapReadString() (string, error) {
	return "", ErrUnsupportedType
}

func (i *FromUrl) MapReadByteString(value []byte) ([]byte, error) {
	return value, ErrUnsupportedType
}

func (i *FromUrl) MapReadBytes(value []byte) ([]byte, error) {
	return value, ErrUnsupportedType
}

func (i *FromUrl) MapReadBool() (bool, error) {
	return false, ErrUnsupportedType
}

func (i *FromUrl) MapReadBigInt(value *big.Int) (*big.Int, error) {

	return value, ErrUnsupportedType
}

func (i *FromUrl) MapReadRat(value *big.Rat) (*big.Rat, error) {

	return value, ErrUnsupportedType
}

func (i *FromUrl) MapReadBigFloat(value *big.Float) (*big.Float, error) {
	return value, ErrUnsupportedType
}

func (i *FromUrl) MapReadValue() error {
	return ErrUnsupportedType
}
