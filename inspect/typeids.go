package inspect

import ()

type TypeId int

const TypeIdName = "inspect.typeid"

const (
	TypeInvalid TypeId = iota
	TypeInt
	TypeInt32
	TypeInt64
	TypeFloat32
	TypeFloat64
	TypeString
	TypeByteString
	TypeBytes
	TypeBool
	TypeBigInt
	TypeRat
	TypeBigFloat
	TypeObject
	TypeArray
	TypeValue
	TypeMap
	TypeLast
)

func (id *TypeId) Inspect(inspector *GenericInspector) {
	inspector.Int((*int)(id), TypeIdName, "type id for types natively supported by inspectors")
}
