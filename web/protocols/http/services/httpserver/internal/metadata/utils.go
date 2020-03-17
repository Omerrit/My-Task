package metadata

import (
	"gerrit-share.lan/go/inspect"
	"log"
)

func IsNested(data *Metadata) bool {
	for _, underlyingValue := range data.UnderlyingValues {
		if underlyingValue.TypeId == inspect.TypeValue {
			for ; underlyingValue.TypeId != inspect.TypeValue; underlyingValue = underlyingValue.UnderlyingValues[0] {
			}
		}
		if underlyingValue.UnderlyingValues != nil {
			return true
		}
	}
	return false
}

func MakeCommandMetaData(command inspect.Inspectable) (*Metadata, error) {
	inspectorImpl := NewMetadataCreator()
	inspector := inspect.NewGenericInspector(inspectorImpl)
	command.Inspect(inspector)
	return inspectorImpl.Metadata, inspector.GetError()
}

var typeNames = []string{
	"", "int", "int32", "int64", "float32", "float64", "string", "string", "bytes", "bool",
	"big int", "big rat", "big float", "object", "array", "value wrapper", "map",
}

func GetTypeName(id inspect.TypeId) string {
	if id >= inspect.TypeLast {
		log.Printf("unsupported type id has been detected: %v", id)
		return ""
	}
	return typeNames[id]
}

func GetTypeIdByName(typeName string) inspect.TypeId {
	for id, val := range typeNames {
		if val == typeName {
			return inspect.TypeId(id)
		}
	}
	return inspect.TypeInvalid
}
