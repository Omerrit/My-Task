package arrays

import "gerrit-share.lan/go/inspect"

type StringArray []string

const StringArrayName = packageName + ".string"

func (s *StringArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(StringArrayName, "string", "readable/writable string array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*s))
		} else {
			s.SetLength(arrayInspector.GetLength())
		}
		for index := range *s {
			arrayInspector.String(&(*s)[index])
		}
		arrayInspector.End()
	}
}

func (s *StringArray) SetLength(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([]string, length)
	}
}

func (s *StringArray) Resize(length int) {
	if cap(*s) > length {
		oldLength := len(*s)
		*s = (*s)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*s)[i] = ""
			}
		}
	} else {
		tempSlice := make([]string, length)
		copy(tempSlice, *s)
		*s = tempSlice
	}
}

func (s *StringArray) Push(item string) {
	*s = append(*s, item)
}

func (s *StringArray) Pop() string {
	removed := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return removed
}

func (s *StringArray) IsEmpty() bool {
	return len(*s) == 0
}

func (s *StringArray) Clear() {
	*s = (*s)[:0]
}

type ByteStringArray [][]byte

const ByteStringArrayName = packageName + ".bytestr"

func (b *ByteStringArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(ByteStringArrayName, "byteString", "readable/writable byteString array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*b))
		} else {
			b.SetLength(arrayInspector.GetLength())
		}
		for index := range *b {
			arrayInspector.ByteString(&(*b)[index])
		}
		arrayInspector.End()
	}
}

func (b *ByteStringArray) SetLength(length int) {
	if cap(*b) > length {
		*b = (*b)[:length]
	} else {
		*b = make([][]byte, length)
	}
}

func (b *ByteStringArray) Resize(length int) {
	if cap(*b) > length {
		oldLength := len(*b)
		*b = (*b)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*b)[i] = (*b)[i][:0]
			}
		}
	} else {
		tempSlice := make([][]byte, length)
		copy(tempSlice, *b)
		*b = tempSlice
	}
}

func (b *ByteStringArray) Push(item []byte) {
	*b = append(*b, item)
}

func (b *ByteStringArray) Pop() []byte {
	removed := (*b)[len(*b)-1]
	*b = (*b)[:len(*b)-1]
	return removed
}

func (b *ByteStringArray) IsEmpty() bool {
	return len(*b) == 0
}

func (b *ByteStringArray) Clear() {
	*b = (*b)[:0]
}
