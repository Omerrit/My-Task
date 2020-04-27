package utils

import (
	"fmt"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/json/fromjson"
)

const valueLength = 4

type Value struct {
	//[userName, Date.now(), newValue, message]
	userName string
	time     int
	newValue string
	message  string
}

func (v *Value) UserName() string {
	return v.userName
}

func (v *Value) Time() int {
	return v.time
}

func (v *Value) NewValue() string {
	return v.newValue
}

func (v *Value) Message() string {
	return v.message
}

func (v *Value) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array("value", "", "")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(valueLength)
		} else {
			if arrayInspector.GetLength() != valueLength {
				arrayInspector.SetError(fmt.Errorf("incorrect value format"))
				return
			}
		}
		arrayInspector.String(&v.userName)
		arrayInspector.Int(&v.time)
		arrayInspector.String(&v.newValue)
		arrayInspector.String(&v.message)
		arrayInspector.End()
	}
}

func ParseValue(data []byte) (*Value, error) {
	parser := fromjson.NewInspector(data, 0)
	inspector := inspect.NewGenericInspector(parser)
	value := new(Value)
	value.Inspect(inspector)
	return value, inspector.GetError()
}
