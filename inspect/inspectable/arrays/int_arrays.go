package arrays

import "gerrit-share.lan/go/inspect"

const (
	IntArrayName   = packageName + ".Int"
	Int32ArrayName = packageName + ".Int32"
	Int64ArrayName = packageName + ".Int64"
)

type IntArray []int

func (s *IntArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(IntArrayName, "int", "readable/writable int array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.Int(&(*s)[index])
	}
	arrayInspector.End()
}

func (s *IntArray) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([]int, length)
	}
}

type Int32Array []int32

func (s *Int32Array) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(Int32ArrayName, "int32", "readable/writable int32 array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.Int32(&(*s)[index])
	}
	arrayInspector.End()
}

func (s *Int32Array) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([]int32, length)
	}
}

type Int64Array []int64

func (s *Int64Array) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(Int64ArrayName, "int64", "readable/writable int64 array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.Int64(&(*s)[index])
	}
	arrayInspector.End()
}

func (s *Int64Array) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([]int64, length)
	}
}
