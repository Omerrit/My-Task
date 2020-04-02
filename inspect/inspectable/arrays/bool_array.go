package arrays

import (
	"gerrit-share.lan/go/inspect"
)

type BoolArray []bool

const BoolArrayName = packageName + ".bool"

func (b *BoolArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(BoolArrayName, "bool", "readable/writable bool array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*b))
		} else {
			b.SetLength(arrayInspector.GetLength())
		}
		for index := range *b {
			arrayInspector.Bool(&(*b)[index])
		}
		arrayInspector.End()
	}
}

func (b *BoolArray) SetLength(length int) {
	if cap(*b) > length {
		*b = (*b)[:length]
	} else {
		*b = make([]bool, length)
	}
}

func (b *BoolArray) Resize(length int) {
	if cap(*b) > length {
		oldLength := len(*b)
		*b = (*b)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*b)[i] = false
			}
		}
	} else {
		tempSlice := make([]bool, length)
		copy(tempSlice, *b)
		*b = tempSlice
	}
}

func (b *BoolArray) Push(item bool) {
	*b = append(*b, item)
}

func (b *BoolArray) Pop() bool {
	removed := (*b)[len(*b)-1]
	*b = (*b)[:len(*b)-1]
	return removed
}

func (b *BoolArray) IsEmpty() bool {
	return len(*b) == 0
}

func (b *BoolArray) Clear() {
	*b = (*b)[:0]
}
