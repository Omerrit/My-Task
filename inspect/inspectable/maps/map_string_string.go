package maps

import "gerrit-share.lan/go/inspect"

const MapStringStringName = packageName + ".strstr"

type MapStringString map[string]string

func (m *MapStringString) Inspect(i *inspect.GenericInspector) {
	mapInspector := i.Map(MapStringStringName, "string", "readable/writable string string map")
	{
		if mapInspector.IsReading() {
			length := mapInspector.GetLength()
			if *m == nil {
				*m = make(MapStringString, length)
			}
			for i := 0; i < length; i++ {
				(*m)[mapInspector.NextKey()] = mapInspector.ReadString()
			}
		} else {
			mapInspector.SetLength(len(*m))
			for key, value := range *m {
				mapInspector.WriteString(key, value)
			}
		}
		mapInspector.End()
	}
}
