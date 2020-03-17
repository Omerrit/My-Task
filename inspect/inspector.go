package inspect

import (
	"math/big"
)

type (
	/*
		Serializers or deserializers should implement this.
		GenericInspector guarantees that it wouldn't call any function that can return error
		after some function returned non-nil error

		Each Start* call will have its corresponding End* call

		*Len() functions should use its parameter if inspector is writing (return value doesn't matter)
		and return container length if it's reading (parameter value doesn't matter in this case)

		The value returned by IsReading should never change
		(inspector can be either serializing or deserializing, not both)
	*/
	InspectorImpl interface {
		ValueInt(value *int, typeName string, typeDescription string) error
		ValueInt32(value *int32, typeName string, typeDescription string) error
		ValueInt64(value *int64, typeName string, typeDescription string) error
		ValueFloat32(value *float32, format byte, precision int, typeName string, typeDescription string) error
		ValueFloat64(value *float64, format byte, precision int, typeName string, typeDescription string) error
		ValueString(value *string, typeName string, typeDescription string) error
		ValueByteString(value *[]byte, typeName string, typeDescription string) error
		ValueBytes(value *[]byte, typeName string, typeDescription string) error
		ValueBool(value *bool, typeName string, typeDescription string) error
		ValueBigInt(value *big.Int, typeName string, typeDescription string) error
		ValueRat(value *big.Rat, precision int, typeName string, typeDescription string) error
		ValueBigFloat(value *big.Float, format byte, precision int, typeName string, typeDescription string) error
		StartObject(typeName string, typeDescription string) error
		StartArray(typeName string, valueTypeName string, typeDescription string) error
		StartMap(typeName string, valueTypeName string, typeDescription string) error
		EndObject() error
		EndArray() error
		EndMap() error

		ObjectInt(value *int, name string, mandatory bool, description string) error
		ObjectInt32(value *int32, name string, mandatory bool, description string) error
		ObjectInt64(value *int64, name string, mandatory bool, description string) error
		ObjectFloat32(value *float32, format byte, precision int, name string, mandatory bool, description string) error
		ObjectFloat64(value *float64, format byte, precision int, name string, mandatory bool, description string) error
		ObjectString(value *string, name string, mandatory bool, description string) error
		ObjectByteString(value *[]byte, name string, mandatory bool, description string) error
		ObjectBytes(value *[]byte, name string, mandatory bool, description string) error
		ObjectBool(value *bool, name string, mandatory bool, description string) error
		ObjectBigInt(value *big.Int, name string, mandatory bool, description string) error
		ObjectRat(value *big.Rat, precision int, name string, mandatory bool, description string) error
		ObjectBigFloat(value *big.Float, format byte, precision int, name string, mandatory bool, description string) error
		ObjectValue(name string, mandatory bool, description string) (bool, error)

		ArrayLen(length int) (int, error)
		ArrayInt(value *int) error
		ArrayInt32(value *int32) error
		ArrayInt64(value *int64) error
		ArrayFloat32(value *float32, format byte, precision int) error
		ArrayFloat64(value *float64, format byte, precision int) error
		ArrayString(value *string) error
		ArrayByteString(value *[]byte) error
		ArrayBytes(value *[]byte) error
		ArrayBool(value *bool) error
		ArrayBigInt(value *big.Int) error
		ArrayRat(value *big.Rat, precision int) error
		ArrayBigFloat(value *big.Float, format byte, precision int) error
		ArrayValue() error

		MapLen(length int) (int, error)
		MapNextKey() (string, error)
		MapReadInt() (int, error)
		MapReadInt32() (int32, error)
		MapReadInt64() (int64, error)
		MapReadFloat32() (float32, error)
		MapReadFloat64() (float64, error)
		MapReadString() (string, error)
		MapReadByteString(value []byte) ([]byte, error)
		MapReadBytes(value []byte) ([]byte, error)
		MapReadBool() (bool, error)
		MapReadBigInt(value *big.Int) (*big.Int, error)
		MapReadRat(value *big.Rat) (*big.Rat, error)
		MapReadBigFloat(value *big.Float) (*big.Float, error)
		MapReadValue() error

		MapWriteInt(key string, value int) error
		MapWriteInt32(key string, value int32) error
		MapWriteInt64(key string, value int64) error
		MapWriteFloat32(key string, value float32, format byte, precision int) error
		MapWriteFloat64(key string, value float64, format byte, precision int) error
		MapWriteString(key string, value string) error
		MapWriteByteString(key string, value []byte) error
		MapWriteBytes(key string, value []byte) error
		MapWriteBool(key string, value bool) error
		MapWriteBigInt(key string, value *big.Int) error
		MapWriteRat(key string, value *big.Rat, precision int) error
		MapWriteBigFloat(key string, value *big.Float, format byte, precision int) error
		MapWriteValue(key string) error
		IsReading() bool
	}

	//Type should implement this to be serializable and/or deserializable via GenericInspector
	Inspectable interface {
		Inspect(*GenericInspector)
	}
)

//Empty type
type DummyInspectable struct{}

func (DummyInspectable) Inspect(*GenericInspector) {}

//A type that presents itself as an empty object
type EmptyObject struct{}

func (EmptyObject) Inspect(inspector *GenericInspector) {
	inspector.Object("", "").End()
}

// Fast check if a type implements InspectorImpl
func TestInspectorImpl(InspectorImpl) {}

// Fast check if a type implements Inspectable
func TestInspectable(Inspectable) {}
