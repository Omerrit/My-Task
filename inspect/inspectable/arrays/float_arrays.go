package arrays

import "gerrit-share.lan/go/inspect"

const (
	Float32ArrayName = packageName + ".float32"
	Float64ArrayName = packageName + ".float64"
)

type Float32Array []float32

func (s *Float32Array) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(Float32ArrayName, "float32", "readable/writable float32 array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.Float32(&(*s)[index], 'g', -1)
	}
	arrayInspector.End()
}

func (s *Float32Array) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([]float32, length)
	}
}

type Float64Array []float64

func (s *Float64Array) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(Float64ArrayName, "float64", "readable/writable float64 array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.Float64(&(*s)[index], 'g', -1)
	}
	arrayInspector.End()
}

func (s *Float64Array) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([]float64, length)
	}
}
