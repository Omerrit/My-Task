package inspect

import (
	"math/big"
)

type GenericInspector struct {
	impl       InspectorImpl
	isReading  bool //cached Impl.IsReading
	err        error
	extraValue interface{}
}

func NewGenericInspector(impl InspectorImpl) *GenericInspector {
	return &GenericInspector{impl, impl.IsReading(), nil, nil}
}

func NewGenericInspectorWithExtra(impl InspectorImpl, extraValue interface{}) *GenericInspector {
	return &GenericInspector{impl, impl.IsReading(), nil, extraValue}
}

func (g *GenericInspector) GetExtraValue() interface{} {
	return g.extraValue
}

func (g *GenericInspector) IsReading() bool {
	return g.isReading
}

func (g *GenericInspector) SetError(err error) {
	if err == nil {
		return
	}
	g.err = err
}

func (g *GenericInspector) GetError() error {
	return g.err
}

//value inspector methods

//TODO: implement optional check that value's Inspect can provide only one value

func (g *GenericInspector) Int(value *int, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueInt(value, typeName, typeDescription)
}

func (g *GenericInspector) Int32(value *int32, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueInt32(value, typeName, typeDescription)
}

func (g *GenericInspector) Int64(value *int64, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueInt64(value, typeName, typeDescription)
}

func (g *GenericInspector) Float32(value *float32, format byte, precision int, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueFloat32(value, format, precision, typeName, typeDescription)
}

func (g *GenericInspector) Float64(value *float64, format byte, precision int, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueFloat64(value, format, precision, typeName, typeDescription)
}

func (g *GenericInspector) String(value *string, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueString(value, typeName, typeDescription)
}

func (g *GenericInspector) ByteString(value *[]byte, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueByteString(value, typeName, typeDescription)
}

func (g *GenericInspector) Bytes(value *[]byte, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueBytes(value, typeName, typeDescription)
}

func (g *GenericInspector) Bool(value *bool, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueBool(value, typeName, typeDescription)
}

func (g *GenericInspector) BigInt(value *big.Int, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueBigInt(value, typeName, typeDescription)
}

func (g *GenericInspector) Rat(value *big.Rat, precision int, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueRat(value, precision, typeName, typeDescription)
}

func (g *GenericInspector) BigFloat(value *big.Float, format byte, precision int, typeName string, typeDescription string) {
	if g.err != nil {
		return
	}
	g.err = g.impl.ValueBigFloat(value, format, precision, typeName, typeDescription)
}

func (g *GenericInspector) Object(typeName string, typeDescription string) *ObjectInspector {
	if g.err != nil {
		return (*ObjectInspector)(g)
	}
	g.err = g.impl.StartObject(typeName, typeDescription)
	return (*ObjectInspector)(g)
}

func (g *GenericInspector) Array(typeName string, valueTypeName string, typeDescription string) *ArrayInspector {
	if g.err != nil {
		return (*ArrayInspector)(g)
	}
	g.err = g.impl.StartArray(typeName, valueTypeName, typeDescription)
	return (*ArrayInspector)(g)
}

func (g *GenericInspector) Map(typeName string, valueTypeName string, typeDescription string) *MapInspector {
	if g.err != nil {
		return (*MapInspector)(g)
	}
	g.err = g.impl.StartMap(typeName, valueTypeName, typeDescription)
	return (*MapInspector)(g)
}

type ObjectInspector GenericInspector

func (o *ObjectInspector) IsReading() bool {
	return o.isReading
}

func (o *ObjectInspector) SetError(err error) {
	o.err = err
}

func (o *ObjectInspector) GetExtraValue() interface{} {
	return o.extraValue
}

func (o *ObjectInspector) End() {
	if o.err != nil {
		return
	}
	o.err = o.impl.EndObject()
}

func (o *ObjectInspector) Int(value *int, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectInt(value, name, mandatory, description)
}

func (o *ObjectInspector) Int32(value *int32, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectInt32(value, name, mandatory, description)
}

func (o *ObjectInspector) Int64(value *int64, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectInt64(value, name, mandatory, description)
}

func (o *ObjectInspector) Float32(value *float32, format byte, precision int, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectFloat32(value, format, precision, name, mandatory, description)
}

func (o *ObjectInspector) Float64(value *float64, format byte, precision int, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectFloat64(value, format, precision, name, mandatory, description)
}

func (o *ObjectInspector) String(value *string, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectString(value, name, mandatory, description)
}

func (o *ObjectInspector) ByteString(value *[]byte, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectByteString(value, name, mandatory, description)
}

func (o *ObjectInspector) Bytes(value *[]byte, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectBytes(value, name, mandatory, description)
}

func (o *ObjectInspector) Bool(value *bool, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectBool(value, name, mandatory, description)
}

func (o *ObjectInspector) BigInt(value *big.Int, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectBigInt(value, name, mandatory, description)
}

func (o *ObjectInspector) Rat(value *big.Rat, precision int, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectRat(value, precision, name, mandatory, description)
}

func (o *ObjectInspector) BigFloat(value *big.Float, format byte, precision int, name string, mandatory bool, description string) {
	if o.err != nil {
		return
	}
	o.err = o.impl.ObjectBigFloat(value, format, precision, name, mandatory, description)
}

func (o *ObjectInspector) Value(name string, mandatory bool, description string) *GenericInspector {
	if o.err != nil {
		return (*GenericInspector)(o)
	}
	var present bool
	present, o.err = o.impl.ObjectValue(name, mandatory, description)
	if present {
		return (*GenericInspector)(o)
	} else {
		return nil
	}
}

type ArrayInspector GenericInspector

func (a *ArrayInspector) IsReading() bool {
	return a.isReading
}

func (a *ArrayInspector) SetError(err error) {
	a.err = err
}

func (a *ArrayInspector) GetExtraValue() interface{} {
	return a.extraValue
}

func (a *ArrayInspector) End() {
	if a.err != nil {
		return
	}
	a.err = a.impl.EndArray()
}

func (a *ArrayInspector) GetLength() int {
	if a.err != nil {
		return 0
	}
	var length int
	length, a.err = a.impl.ArrayLen(0)
	return length
}

func (a *ArrayInspector) SetLength(length int) {
	if a.err != nil {
		return
	}
	_, a.err = a.impl.ArrayLen(length)
}

func (a *ArrayInspector) Int(value *int) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayInt(value)
}

func (a *ArrayInspector) Int32(value *int32) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayInt32(value)
}

func (a *ArrayInspector) Int64(value *int64) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayInt64(value)
}

func (a *ArrayInspector) Float32(value *float32, format byte, precision int) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayFloat32(value, format, precision)
}

func (a *ArrayInspector) Float64(value *float64, format byte, precision int) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayFloat64(value, format, precision)
}

func (a *ArrayInspector) String(value *string) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayString(value)
}

func (a *ArrayInspector) ByteString(value *[]byte) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayByteString(value)
}

func (a *ArrayInspector) Bytes(value *[]byte) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayBytes(value)
}

func (a *ArrayInspector) Bool(value *bool) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayBool(value)
}

func (a *ArrayInspector) BigInt(value *big.Int) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayBigInt(value)
}

func (a *ArrayInspector) Rat(value *big.Rat, precision int) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayRat(value, precision)
}

func (a *ArrayInspector) BigFloat(value *big.Float, format byte, precision int) {
	if a.err != nil {
		return
	}
	a.err = a.impl.ArrayBigFloat(value, format, precision)
}

func (a *ArrayInspector) Value() *GenericInspector {
	if a.err != nil {
		return (*GenericInspector)(a)
	}
	a.err = a.impl.ArrayValue()
	return (*GenericInspector)(a)
}

type MapInspector GenericInspector

func (m *MapInspector) IsReading() bool {
	return m.isReading
}

func (m *MapInspector) SetError(err error) {
	m.err = err
}

func (g *MapInspector) GetExtraValue() interface{} {
	return g.extraValue
}

func (m *MapInspector) End() {
	if m.err != nil {
		return
	}
	m.err = m.impl.EndMap()
}

func (m *MapInspector) GetLength() int {
	if m.err != nil {
		return 0
	}
	var length int
	length, m.err = m.impl.MapLen(0)
	return length
}

func (m *MapInspector) SetLength(length int) {
	if m.err != nil {
		return
	}
	_, m.err = m.impl.MapLen(length)
}

//do we need separate functions or separate inspectors for reading/writing?

//reading functions
func (m *MapInspector) NextKey() string {
	if m.err != nil {
		return ""
	}
	var result string
	result, m.err = m.impl.MapNextKey()
	return result
}

func (m *MapInspector) ReadInt() int {
	if m.err != nil {
		return 0
	}
	var result int
	result, m.err = m.impl.MapReadInt()
	return result
}

func (m *MapInspector) ReadInt32() int32 {
	if m.err != nil {
		return 0
	}
	var result int32
	result, m.err = m.impl.MapReadInt32()
	return result
}

func (m *MapInspector) ReadInt64() int64 {
	if m.err != nil {
		return 0
	}
	var result int64
	result, m.err = m.impl.MapReadInt64()
	return result
}

func (m *MapInspector) ReadFloat32() float32 {
	if m.err != nil {
		return 0
	}
	var result float32
	result, m.err = m.impl.MapReadFloat32()
	return result
}

func (m *MapInspector) ReadFloat64() float64 {
	if m.err != nil {
		return 0
	}
	var result float64
	result, m.err = m.impl.MapReadFloat64()
	return result
}

func (m *MapInspector) ReadString() string {
	if m.err != nil {
		return ""
	}
	var result string
	result, m.err = m.impl.MapReadString()
	return result
}

func (m *MapInspector) ReadByteString(value []byte) []byte {
	if m.err != nil {
		return value
	}
	var result []byte
	result, m.err = m.impl.MapReadByteString(value)
	return result
}

func (m *MapInspector) ReadBytes(value []byte) []byte {
	if m.err != nil {
		return value
	}
	var result []byte
	result, m.err = m.impl.MapReadBytes(value)
	return result
}

func (m *MapInspector) ReadBool() bool {
	if m.err != nil {
		return false
	}
	var result bool
	result, m.err = m.impl.MapReadBool()
	return result
}

func (m *MapInspector) ReadBigInt(value *big.Int) *big.Int {
	if m.err != nil {
		return value
	}
	var result *big.Int
	result, m.err = m.impl.MapReadBigInt(value)
	return result
}

func (m *MapInspector) ReadRat(value *big.Rat) *big.Rat {
	if m.err != nil {
		return value
	}
	var result *big.Rat
	result, m.err = m.impl.MapReadRat(value)
	return result
}

func (m *MapInspector) ReadBigFloat(value *big.Float) *big.Float {
	if m.err != nil {
		return value
	}
	var result *big.Float
	result, m.err = m.impl.MapReadBigFloat(value)
	return result
}

func (m *MapInspector) ReadValue() *GenericInspector {
	if m.err != nil {
		return (*GenericInspector)(m)
	}
	m.err = m.impl.MapReadValue()
	return (*GenericInspector)(m)
}

func (m *MapInspector) WriteInt(key string, value int) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteInt(key, value)
}

func (m *MapInspector) WriteInt32(key string, value int32) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteInt32(key, value)
}

func (m *MapInspector) WriteInt64(key string, value int64) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteInt64(key, value)
}

func (m *MapInspector) WriteFloat32(key string, value float32, format byte, precision int) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteFloat32(key, value, format, precision)
}

func (m *MapInspector) WriteFloat64(key string, value float64, format byte, precision int) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteFloat64(key, value, format, precision)
}

func (m *MapInspector) WriteString(key string, value string) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteString(key, value)
}

func (m *MapInspector) WriteByteString(key string, value []byte) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteByteString(key, value)
}

func (m *MapInspector) WriteBytes(key string, value []byte) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteBytes(key, value)
}

func (m *MapInspector) WriteBool(key string, value bool) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteBool(key, value)
}

func (m *MapInspector) WriteBigInt(key string, value *big.Int) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteBigInt(key, value)
}

func (m *MapInspector) WriteRat(key string, value *big.Rat, precision int) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteRat(key, value, precision)
}

func (m *MapInspector) WriteBigFloat(key string, value *big.Float, format byte, precision int) {
	if m.err != nil {
		return
	}
	m.err = m.impl.MapWriteBigFloat(key, value, format, precision)
}

func (m *MapInspector) WriteValue(key string) *GenericInspector {
	if m.err != nil {
		return (*GenericInspector)(m)
	}
	m.err = m.impl.MapWriteValue(key)
	return (*GenericInspector)(m)
}
