package arrays

import (
	"gerrit-share.lan/go/inspect"
)

type IntArray []int

const IntArrayName = packageName + ".int"

func (a *IntArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(IntArrayName, "int", "readable/writable int array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*a))
		} else {
			a.SetLength(arrayInspector.GetLength())
		}
		for index := range *a {
			arrayInspector.Int(&(*a)[index])
		}
		arrayInspector.End()
	}
}

func (a *IntArray) SetLength(length int) {
	if cap(*a) > length {
		*a = (*a)[:length]
	} else {
		*a = make([]int, length)
	}
}

func (a *IntArray) Resize(length int) {
	if cap(*a) > length {
		oldLength := len(*a)
		*a = (*a)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*a)[i] = 0
			}
		}
	} else {
		tempSlice := make([]int, length)
		copy(tempSlice, *a)
		*a = tempSlice
	}
}

func (a *IntArray) Push(item int) {
	*a = append(*a, item)
}

func (a *IntArray) Pop() int {
	removed := (*a)[len(*a)-1]
	*a = (*a)[:len(*a)-1]
	return removed
}

func (a *IntArray) IsEmpty() bool {
	return len(*a) == 0
}

func (a *IntArray) Clear() {
	*a = (*a)[:0]
}

type Int32Array []int32

const Int32ArrayName = packageName + ".int32"

func (a *Int32Array) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(Int32ArrayName, "int32", "readable/writable int32 array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*a))
		} else {
			a.SetLength(arrayInspector.GetLength())
		}
		for index := range *a {
			arrayInspector.Int32(&(*a)[index])
		}
		arrayInspector.End()
	}
}

func (a *Int32Array) SetLength(length int) {
	if cap(*a) > length {
		*a = (*a)[:length]
	} else {
		*a = make([]int32, length)
	}
}

func (a *Int32Array) Resize(length int) {
	if cap(*a) > length {
		oldLength := len(*a)
		*a = (*a)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*a)[i] = 0
			}
		}
	} else {
		tempSlice := make([]int32, length)
		copy(tempSlice, *a)
		*a = tempSlice
	}
}

func (a *Int32Array) Push(item int32) {
	*a = append(*a, item)
}

func (a *Int32Array) Pop() int32 {
	removed := (*a)[len(*a)-1]
	*a = (*a)[:len(*a)-1]
	return removed
}

func (a *Int32Array) IsEmpty() bool {
	return len(*a) == 0
}

func (a *Int32Array) Clear() {
	*a = (*a)[:0]
}

type Int64Array []int64

const Int64ArrayName = packageName + ".int64"

func (a *Int64Array) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(Int64ArrayName, "int64", "readable/writable int64 array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*a))
		} else {
			a.SetLength(arrayInspector.GetLength())
		}
		for index := range *a {
			arrayInspector.Int64(&(*a)[index])
		}
		arrayInspector.End()
	}
}

func (a *Int64Array) SetLength(length int) {
	if cap(*a) > length {
		*a = (*a)[:length]
	} else {
		*a = make([]int64, length)
	}
}

func (a *Int64Array) Resize(length int) {
	if cap(*a) > length {
		oldLength := len(*a)
		*a = (*a)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*a)[i] = 0
			}
		}
	} else {
		tempSlice := make([]int64, length)
		copy(tempSlice, *a)
		*a = tempSlice
	}
}

func (a *Int64Array) Push(item int64) {
	*a = append(*a, item)
}

func (a *Int64Array) Pop() int64 {
	removed := (*a)[len(*a)-1]
	*a = (*a)[:len(*a)-1]
	return removed
}

func (a *Int64Array) IsEmpty() bool {
	return len(*a) == 0
}

func (a *Int64Array) Clear() {
	*a = (*a)[:0]
}
