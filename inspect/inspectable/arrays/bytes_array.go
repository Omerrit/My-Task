package arrays

import (
	"gerrit-share.lan/go/inspect"
)

type BytesArray [][]byte

const BytesArrayName = packageName + ".bytes"

func (b *BytesArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(BytesArrayName, "bytes", "readable/writable byte array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*b))
		} else {
			b.SetLength(arrayInspector.GetLength())
		}
		for index := range *b {
			arrayInspector.Bytes(&(*b)[index])
		}
		arrayInspector.End()
	}
}

func (b *BytesArray) SetLength(length int) {
	if cap(*b) > length {
		*b = (*b)[:length]
	} else {
		*b = make([][]byte, length)
	}
}

func (b *BytesArray) Resize(length int) {
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

func (b *BytesArray) Push(item []byte) {
	*b = append(*b, item)
}

func (b *BytesArray) Pop() []byte {
	removed := (*b)[len(*b)-1]
	*b = (*b)[:len(*b)-1]
	return removed
}

func (b *BytesArray) IsEmpty() bool {
	return len(*b) == 0
}

func (b *BytesArray) Clear() {
	*b = (*b)[:0]
}
