package frombytes

import (
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect/json/fromjson"
	"math/big"
)

type ProxyPathValuesInspector struct {
	ProxyBatchInspector
	path         []byte
	isValues     bool
	acceptValues bool
}

func NewProxyInspector(path []byte, data []byte, offset int, acceptValues bool) *ProxyPathValuesInspector {
	return &ProxyPathValuesInspector{ProxyBatchInspector: *NewProxyBatchInspector(data, offset), path: path, acceptValues: acceptValues}
}

func (i *ProxyPathValuesInspector) StartObject(typeName string, typeDescription string) error {
	err := i.InspectorImpl.StartObject(typeName, typeDescription)
	i.depth++
	if i.depth > 2 {
		return err
	}
	if err != nil {
		if !i.acceptValues {
			return ErrTooFewParameters
		}
		i.isValues = true
	}
	return nil
}

func (i *ProxyPathValuesInspector) EndObject() error {
	i.depth--
	if i.depth == 1 && i.isValues {
		return nil
	}
	return i.InspectorImpl.EndObject()
}

func (i *ProxyPathValuesInspector) ObjectInt(value *int, name string, mandatory bool, description string) error {
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueInt(value, "", "")
	}
	return i.InspectorImpl.ObjectInt(value, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectString(value *string, name string, mandatory bool, description string) error {
	if name == "path" {
		err := i.InspectorImpl.ObjectString(value, name, true, description)
		if err != nil {
			if errors.Is(err, fromjson.ErrPropertyNotFound) {
				if mandatory && len(i.path) == 0 {
					return err
				}
				*value = string(i.path)
				return nil
			}
			return err
		}
		if len(i.path) > 0 {
			return ErrDuplicatePath
		}
		return nil
	}
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueString(value, "", "")
	}
	return i.InspectorImpl.ObjectString(value, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectByteString(value *[]byte, name string, mandatory bool, description string) error {
	if name == "path" {
		err := i.InspectorImpl.ObjectByteString(value, name, true, description)
		if err != nil {
			if errors.Is(err, fromjson.ErrPropertyNotFound) {
				if mandatory && len(i.path) == 0 {
					return err
				}
				*value = i.path
				return nil
			}
			return err
		}
		if len(i.path) > 0 {
			return ErrDuplicatePath
		}
		return nil
	}
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueByteString(value, "", "")
	}
	return i.InspectorImpl.ObjectByteString(value, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectInt32(value *int32, name string, mandatory bool, description string) error {
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueInt32(value, "", "")
	}
	return i.InspectorImpl.ObjectInt32(value, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectInt64(value *int64, name string, mandatory bool, description string) error {
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueInt64(value, "", "")
	}
	return i.InspectorImpl.ObjectInt64(value, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectFloat32(value *float32, format byte, precision int, name string, mandatory bool, description string) error {
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueFloat32(value, format, precision, "", "")
	}
	return i.InspectorImpl.ObjectFloat32(value, format, precision, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectFloat64(value *float64, format byte, precision int, name string, mandatory bool, description string) error {
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueFloat64(value, format, precision, "", "")
	}
	return i.InspectorImpl.ObjectFloat64(value, format, precision, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectBytes(value *[]byte, name string, mandatory bool, description string) error {
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueBytes(value, "", "")
	}
	return i.InspectorImpl.ObjectBytes(value, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectBool(value *bool, name string, mandatory bool, description string) error {
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueBool(value, "", "")
	}
	return i.InspectorImpl.ObjectBool(value, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectBigInt(value *big.Int, name string, mandatory bool, description string) error {
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueBigInt(value, "", "")
	}
	return i.InspectorImpl.ObjectBigInt(value, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectRat(value *big.Rat, precision int, name string, mandatory bool, description string) error {
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueRat(value, precision, "", "")
	}
	return i.InspectorImpl.ObjectRat(value, precision, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectBigFloat(value *big.Float, format byte, precision int, name string, mandatory bool, description string) error {
	if i.depth == 2 && i.isValues {
		return i.InspectorImpl.ValueBigFloat(value, format, precision, "", "")
	}
	return i.InspectorImpl.ObjectBigFloat(value, format, precision, name, mandatory, description)
}

func (i *ProxyPathValuesInspector) ObjectValue(name string, mandatory bool, description string) (bool, error) {
	if i.depth == 2 && i.isValues {
		return true, nil
	}
	return i.InspectorImpl.ObjectValue(name, mandatory, description)
}
