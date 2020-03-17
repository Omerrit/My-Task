package frombytes

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/json/fromjson"
)

type ProxyBatchInspector struct {
	inspect.InspectorImpl
	IsBatch bool
	depth   int
}

func NewProxyBatchInspector(data []byte, offset int) *ProxyBatchInspector {
	return &ProxyBatchInspector{InspectorImpl: fromjson.NewInspector(data, offset)}
}

func (i *ProxyBatchInspector) StartArray(typeName string, valueTypeName string, typeDescription string) error {
	err := i.InspectorImpl.StartArray(typeName, valueTypeName, typeDescription)
	i.depth++
	if i.depth > 1 {
		return err
	}
	if err == nil {
		i.IsBatch = true
	}
	return nil
}

func (i *ProxyBatchInspector) EndArray() error {
	i.depth--
	if i.depth == 0 && !i.IsBatch {
		return nil
	}
	return i.InspectorImpl.EndArray()
}

func (i *ProxyBatchInspector) StartObject(typeName string, typeDescription string) error {
	i.depth++
	return i.InspectorImpl.StartObject(typeName, typeDescription)
}

func (i *ProxyBatchInspector) EndObject() error {
	i.depth--
	return i.InspectorImpl.EndObject()
}

func (i *ProxyBatchInspector) StartMap(typeName string, valueTypeName string, typeDescription string) error {
	i.depth++
	return i.InspectorImpl.StartMap(typeName, valueTypeName, typeDescription)
}

func (i *ProxyBatchInspector) EndMap() error {
	i.depth--
	return i.InspectorImpl.EndMap()
}

func (i *ProxyBatchInspector) ArrayLen(length int) (int, error) {
	if i.IsBatch {
		return i.InspectorImpl.ArrayLen(length)
	}
	return 1, nil
}

func (i *ProxyBatchInspector) ArrayValue() error {
	if i.depth == 1 && !i.IsBatch {
		return nil
	}
	return i.InspectorImpl.ArrayValue()
}
