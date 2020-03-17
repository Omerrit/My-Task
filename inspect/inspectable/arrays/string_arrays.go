package arrays

import "gerrit-share.lan/go/inspect"

const (
	StringArrayName     = packageName + ".string"
	ByteStringArrayName = packageName + ".bytestr"
)

type StringArray []string

func (s *StringArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(StringArrayName, "string", "readable/writable string array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.String(&(*s)[index])
	}
	arrayInspector.End()
}

func (s *StringArray) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([]string, length)
	}
}

type ByteStringArray [][]byte

func (s *ByteStringArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(ByteStringArrayName, "byteString", "readable/writable byteString array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.ByteString(&(*s)[index])
	}
	arrayInspector.End()
}

func (s *ByteStringArray) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([][]byte, length)
	}
}
