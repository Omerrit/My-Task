package arrays

import "gerrit-share.lan/go/inspect"

const BytesArrayName = packageName + ".bytes"

type BytesArray [][]byte

func (s *BytesArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(BytesArrayName, "bytes", "readable/writable byte array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.Bytes(&(*s)[index])
	}
	arrayInspector.End()
}

func (s *BytesArray) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([][]byte, length)
	}
}
