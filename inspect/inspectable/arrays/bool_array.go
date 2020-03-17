package arrays

import "gerrit-share.lan/go/inspect"

const BoolArrayName = packageName + ".bool"

type BoolArray []bool

func (b *BoolArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(BoolArrayName, "bool", "readable/writable bool array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*b))
	} else {
		b.init(arrayInspector.GetLength())
	}
	for index := range *b {
		arrayInspector.Bool(&(*b)[index])
	}
	arrayInspector.End()
}

func (b *BoolArray) init(length int) {
	if cap(*b) > length {
		*b = (*b)[:length]
	} else {
		*b = make([]bool, length)
	}
}
