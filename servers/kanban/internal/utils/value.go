package utils

import (
	"fmt"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/json/fromjson"
	"gerrit-share.lan/go/inspect/json/tojson"
	"time"
)

const valueLength = 4

type Value struct {
	//[userName, Date.now(), newValue, message]
	userName string
	time     int64
	newValue string
	message  string
}

func NewValue(userName string, time int64, newValue string, message string) *Value {
	return &Value{userName, time, newValue, message}
}

func (v *Value) ToJson() string {
	serializer := &tojson.Inspector{}
	inspector := inspect.NewGenericInspector(serializer)
	v.Inspect(inspector)
	return string(serializer.Output())
}

func (v *Value) UserName() string {
	return v.userName
}

func (v *Value) Time() time.Time {
	return time.Unix(v.time, 0)
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
		arrayInspector.Int64(&v.time)
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
