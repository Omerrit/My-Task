package metadata

import (
	"encoding/base64"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"gerrit-share.lan/go/inspect/inspectorhelpers"
	"math/big"
	"strconv"
)

type MetadataCreator struct {
	inspectorhelpers.Writer
	Metadata      *Metadata
	prevMetadata  []*Metadata
	currentIndex  int
	typeNames     []string
	ignoreCounter int
}

func newMetadataCreator(valueNames []string) *MetadataCreator {
	data := NewMetadata(inspect.TypeInvalid, "", "", false, "")
	return &MetadataCreator{
		Metadata:     data,
		prevMetadata: []*Metadata{data},
		typeNames:    valueNames,
	}
}

func NewMetadataCreator() *MetadataCreator {
	return newMetadataCreator(nil)
}

func (m *MetadataCreator) ValueInt(value *int, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeInt, typeName, typeDescription, strconv.Itoa(*value))
	return m.end()
}

func (m *MetadataCreator) ValueInt32(value *int32, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeInt32, typeName, typeDescription, strconv.FormatInt(int64(*value), 10))
	return m.end()
}

func (m *MetadataCreator) ValueInt64(value *int64, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeInt64, typeName, typeDescription, strconv.FormatInt(*value, 10))
	return m.end()
}

func (m *MetadataCreator) ValueFloat32(value *float32, format byte, precision int, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeFloat32, typeName, typeDescription, strconv.FormatFloat(float64(*value), format, precision, 64))
	return m.end()
}

func (m *MetadataCreator) ValueFloat64(value *float64, format byte, precision int, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeFloat64, typeName, typeDescription, strconv.FormatFloat(*value, format, precision, 64))
	return m.end()
}

func (m *MetadataCreator) ValueString(value *string, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeString, typeName, typeDescription, *value)
	return m.end()
}

func (m *MetadataCreator) ValueByteString(value *[]byte, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeByteString, typeName, typeDescription, string(*value))
	return m.end()
}

func (m *MetadataCreator) ValueBytes(value *[]byte, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeBytes, typeName, typeDescription, base64.StdEncoding.EncodeToString(*value))
	return m.end()
}

func (m *MetadataCreator) ValueBool(value *bool, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeBool, typeName, typeDescription, strconv.FormatBool(*value))
	return m.end()
}

func (m *MetadataCreator) ValueBigInt(value *big.Int, typeName string, typeDescription string) error {
	defaultVal, _ := value.MarshalText()
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeBigInt, typeName, typeDescription, string(defaultVal))
	return m.end()
}

func (m *MetadataCreator) ValueRat(value *big.Rat, precision int, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeRat, typeName, typeDescription, value.FloatString(precision))
	return m.end()
}

func (m *MetadataCreator) ValueBigFloat(value *big.Float, format byte, precision int, typeName string, typeDescription string) error {
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(inspect.TypeBigFloat, typeName, typeDescription, value.Text(format, precision))
	return m.end()
}

func (m *MetadataCreator) ignoreIfExists(typeName string) bool {
	for _, v := range m.typeNames {
		if v == typeName {
			m.ignoreCounter++
			return true
		}
	}
	return false
}

func (m *MetadataCreator) startValue(typeId inspect.TypeId, typeName, typeDescription, valueTypeName string) error {
	if m.ignoreCounter > 0 {
		m.ignoreCounter++
		return nil
	}
	m.prevMetadata[len(m.prevMetadata)-1].SetTypeInfo(typeId, typeName, typeDescription, "null")
	m.prevMetadata[len(m.prevMetadata)-1].ValueTypeName = valueTypeName
	if m.ignoreIfExists(typeName) {
		return nil
	}
	m.typeNames = append(m.typeNames, typeName)
	return nil
}

func (m *MetadataCreator) StartObject(typeName string, typeDescription string) error {
	return m.startValue(inspect.TypeObject, typeName, typeDescription, "")
}

func (m *MetadataCreator) StartArray(typeName string, valueTypeName string, typeDescription string) error {
	return m.startValue(inspect.TypeArray, typeName, typeDescription, valueTypeName)
}

func (m *MetadataCreator) StartMap(typeName string, valueTypeName string, typeDescription string) error {
	return m.startValue(inspect.TypeMap, typeName, typeDescription, valueTypeName)
}

func (m *MetadataCreator) end() error {
	m.prevMetadata = m.prevMetadata[:len(m.prevMetadata)-1]
	return nil
}

func (m *MetadataCreator) EndObject() error {
	var changed bool
	if m.ignoreCounter > 0 {
		changed = true
		m.ignoreCounter--
	}
	if m.ignoreCounter > 0 {
		return nil
	}
	if !changed {
		m.typeNames = m.typeNames[:len(m.typeNames)-1]
	}
	return m.end()
}

func (m *MetadataCreator) EndArray() error {
	var changed bool
	if m.ignoreCounter > 0 {
		changed = true
		m.ignoreCounter--
	}
	if m.ignoreCounter > 0 {
		return nil
	}
	data := m.prevMetadata[len(m.prevMetadata)-1]
	var duplicateType bool
	for _, v := range m.typeNames {
		if v == data.ValueTypeName {
			duplicateType = true
		}
	}
	if len(data.UnderlyingValues) == 0 {
		creator := inspectables.Get(data.ValueTypeName)
		if creator != nil && !duplicateType {
			reader := newMetadataCreator(m.typeNames)
			serializer := inspect.NewGenericInspector(reader)
			creator().Inspect(serializer)
			data.Add(reader.Metadata)
		} else {
			typeId := GetTypeIdByName(data.ValueTypeName)
			if typeId != inspect.TypeInvalid {
				data.Add(&Metadata{TypeId: typeId})
			}
		}
	}
	if !changed {
		m.typeNames = m.typeNames[:len(m.typeNames)-1]
	}
	return m.end()
}

func (m *MetadataCreator) EndMap() error {
	return m.EndArray()
}

func (m *MetadataCreator) nextIndex() {
	if len(m.prevMetadata) > 1 {
		return
	}
	m.currentIndex++
}

func (m *MetadataCreator) ObjectInt(value *int, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeInt, strconv.Itoa(*value), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectInt32(value *int32, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeInt32, strconv.FormatInt(int64(*value), 10), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectInt64(value *int64, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeInt64, strconv.FormatInt(*value, 10), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectFloat32(value *float32, format byte, precision int, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeFloat32, strconv.FormatFloat(float64(*value), format, precision, 32), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectFloat64(value *float64, format byte, precision int, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeFloat64, strconv.FormatFloat(*value, format, precision, 64), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectString(value *string, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	if len(m.prevMetadata) == 1 && name == "path" {
		m.prevMetadata[0].PathIndex = m.currentIndex
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeString, *value, name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectByteString(value *[]byte, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	if len(m.prevMetadata) == 1 && name == "path" {
		m.prevMetadata[0].PathIndex = m.currentIndex
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeByteString, string(*value), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectBytes(value *[]byte, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeBytes, base64.StdEncoding.EncodeToString(*value), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectBool(value *bool, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeBool, strconv.FormatBool(*value), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectBigInt(value *big.Int, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	m.nextIndex()
	defaultVal, _ := value.MarshalText()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeBigInt, string(defaultVal), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectRat(value *big.Rat, precision int, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeRat, value.FloatString(precision), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectBigFloat(value *big.Float, format byte, precision int, name string, mandatory bool, description string) error {
	if m.ignoreCounter > 0 {
		return nil
	}
	m.nextIndex()
	m.prevMetadata[len(m.prevMetadata)-1].Add(NewMetadata(inspect.TypeBigFloat, value.Text(format, precision), name, mandatory, description))
	return nil
}

func (m *MetadataCreator) ObjectValue(name string, mandatory bool, description string) (bool, error) {
	if m.ignoreCounter > 0 {
		return true, nil
	}
	m.nextIndex()
	data := NewMetadata(inspect.TypeInvalid, "", name, mandatory, description)
	(m.prevMetadata[len(m.prevMetadata)-1]).Add(data)
	m.prevMetadata = append(m.prevMetadata, data)
	return true, nil
}

func (m *MetadataCreator) ArrayLen(length int) (int, error) {
	return 0, nil
}
func (m *MetadataCreator) ArrayInt(value *int) error {
	return m.ObjectInt(value, "", false, "")
}
func (m *MetadataCreator) ArrayInt32(value *int32) error {
	return m.ObjectInt32(value, "", false, "")
}
func (m *MetadataCreator) ArrayInt64(value *int64) error {
	return m.ObjectInt64(value, "", false, "")
}
func (m *MetadataCreator) ArrayFloat32(value *float32, format byte, precision int) error {
	return m.ObjectFloat32(value, format, precision, "", false, "")
}
func (m *MetadataCreator) ArrayFloat64(value *float64, format byte, precision int) error {
	return m.ObjectFloat64(value, format, precision, "", false, "")
}
func (m *MetadataCreator) ArrayString(value *string) error {
	return m.ObjectString(value, "", false, "")
}
func (m *MetadataCreator) ArrayByteString(value *[]byte) error {
	return m.ObjectByteString(value, "", false, "")
}
func (m *MetadataCreator) ArrayBytes(value *[]byte) error {
	return m.ObjectBytes(value, "", false, "")
}
func (m *MetadataCreator) ArrayBool(value *bool) error {
	return m.ObjectBool(value, "", false, "")
}
func (m *MetadataCreator) ArrayBigInt(value *big.Int) error {
	return m.ObjectBigInt(value, "", false, "")
}
func (m *MetadataCreator) ArrayRat(value *big.Rat, precision int) error {
	return m.ObjectRat(value, precision, "", false, "")
}
func (m *MetadataCreator) ArrayBigFloat(value *big.Float, format byte, precision int) error {
	return m.ObjectBigFloat(value, format, precision, "", false, "")
}
func (m *MetadataCreator) ArrayValue() error {
	_, err := m.ObjectValue("", false, "")
	return err
}

func (m *MetadataCreator) MapLen(length int) (int, error) {
	return length, nil
}

func (m *MetadataCreator) MapWriteInt(key string, value int) error {
	return m.ObjectInt(&value, key, false, "")
}

func (m *MetadataCreator) MapWriteInt32(key string, value int32) error {
	return m.ObjectInt32(&value, key, false, "")
}

func (m *MetadataCreator) MapWriteInt64(key string, value int64) error {
	return m.ObjectInt64(&value, key, false, "")
}

func (m *MetadataCreator) MapWriteFloat32(key string, value float32, format byte, precision int) error {
	return m.ObjectFloat32(&value, format, precision, key, false, "")
}

func (m *MetadataCreator) MapWriteFloat64(key string, value float64, format byte, precision int) error {
	return m.ObjectFloat64(&value, format, precision, key, false, "")
}

func (m *MetadataCreator) MapWriteString(key string, value string) error {
	return m.ObjectString(&value, key, false, "")
}

func (m *MetadataCreator) MapWriteByteString(key string, value []byte) error {
	return m.ObjectByteString(&value, key, false, "")
}

func (m *MetadataCreator) MapWriteBytes(key string, value []byte) error {
	return m.ObjectBytes(&value, key, false, "")
}

func (m *MetadataCreator) MapWriteBool(key string, value bool) error {
	return m.ObjectBool(&value, key, false, "")
}

func (m *MetadataCreator) MapWriteBigInt(key string, value *big.Int) error {
	return m.ObjectBigInt(value, key, false, "")
}

func (m *MetadataCreator) MapWriteRat(key string, value *big.Rat, precision int) error {
	return m.ObjectRat(value, precision, key, false, "")
}

func (m *MetadataCreator) MapWriteBigFloat(key string, value *big.Float, format byte, precision int) error {
	return m.ObjectBigFloat(value, format, precision, key, false, "")
}

func (m *MetadataCreator) MapWriteValue(key string) error {
	_, err := m.ObjectValue(key, false, "")
	return err
}

func (m *MetadataCreator) IsReading() bool {
	return false
}
