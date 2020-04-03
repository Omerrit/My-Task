package arrays

import (
	"gerrit-share.lan/go/inspect"
)

type Float32Array []float32

const Float32ArrayName = packageName + ".float32"

func (f *Float32Array) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(Float32ArrayName, "float32", "readable/writable float32 array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*f))
		} else {
			f.SetLength(arrayInspector.GetLength())
		}
		for index := range *f {
			arrayInspector.Float32(&(*f)[index], 'g', -1)
		}
		arrayInspector.End()
	}
}

func (f *Float32Array) SetLength(length int) {
	if cap(*f) > length {
		*f = (*f)[:length]
	} else {
		*f = make([]float32, length)
	}
}

func (f *Float32Array) Resize(length int) {
	if cap(*f) > length {
		oldLength := len(*f)
		*f = (*f)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*f)[i] = 0
			}
		}
	} else {
		tempSlice := make([]float32, length)
		copy(tempSlice, *f)
		*f = tempSlice
	}
}

func (f *Float32Array) Push(item float32) {
	*f = append(*f, item)
}

func (f *Float32Array) Pop() float32 {
	removed := (*f)[len(*f)-1]
	*f = (*f)[:len(*f)-1]
	return removed
}

func (f *Float32Array) IsEmpty() bool {
	return len(*f) == 0
}

func (f *Float32Array) Clear() {
	*f = (*f)[:0]
}

type Float64Array []float64

const Float64ArrayName = packageName + ".float64"

func (f *Float64Array) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(Float64ArrayName, "float64", "readable/writable float64 array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*f))
		} else {
			f.SetLength(arrayInspector.GetLength())
		}
		for index := range *f {
			arrayInspector.Float64(&(*f)[index], 'g', -1)
		}
		arrayInspector.End()
	}
}

func (f *Float64Array) SetLength(length int) {
	if cap(*f) > length {
		*f = (*f)[:length]
	} else {
		*f = make([]float64, length)
	}
}

func (f *Float64Array) Resize(length int) {
	if cap(*f) > length {
		oldLength := len(*f)
		*f = (*f)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*f)[i] = 0
			}
		}
	} else {
		tempSlice := make([]float64, length)
		copy(tempSlice, *f)
		*f = tempSlice
	}
}

func (f *Float64Array) Push(item float64) {
	*f = append(*f, item)
}

func (f *Float64Array) Pop() float64 {
	removed := (*f)[len(*f)-1]
	*f = (*f)[:len(*f)-1]
	return removed
}

func (f *Float64Array) IsEmpty() bool {
	return len(*f) == 0
}

func (f *Float64Array) Clear() {
	*f = (*f)[:0]
}
