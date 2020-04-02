package arrays

import (
	"gerrit-share.lan/go/inspect"
	"math/big"
)

type BigIntArray []big.Int

const BigIntArrayName = packageName + ".bigint"

func (b *BigIntArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(BigIntArrayName, "bigInt", "readable/writable bigInt array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*b))
		} else {
			b.SetLength(arrayInspector.GetLength())
		}
		for index := range *b {
			arrayInspector.BigInt(&(*b)[index])
		}
		arrayInspector.End()
	}
}

func (b *BigIntArray) SetLength(length int) {
	if cap(*b) > length {
		*b = (*b)[:length]
	} else {
		*b = make([]big.Int, length)
	}
}

func (b *BigIntArray) Resize(length int) {
	if cap(*b) > length {
		oldLength := len(*b)
		*b = (*b)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*b)[i].SetInt64(0)
			}
		}
	} else {
		tempSlice := make([]big.Int, length)
		copy(tempSlice, *b)
		*b = tempSlice
	}
}

func (b *BigIntArray) Push(item big.Int) {
	*b = append(*b, item)
}

func (b *BigIntArray) Pop() *big.Int {
	removed := (*b)[len(*b)-1]
	(*b)[len(*b)-1] = big.Int{}
	*b = (*b)[:len(*b)-1]
	return &removed
}

func (b *BigIntArray) IsEmpty() bool {
	return len(*b) == 0
}

func (b *BigIntArray) Clear() {
	*b = (*b)[:0]
}

type BigFloatArray []big.Float

const BigFloatArrayName = packageName + ".bigfloat"

func (b *BigFloatArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(BigFloatArrayName, "bigFloat", "readable/writable bigFloat array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*b))
		} else {
			b.SetLength(arrayInspector.GetLength())
		}
		for index := range *b {
			arrayInspector.BigFloat(&(*b)[index], 'g', -1)
		}
		arrayInspector.End()
	}
}

func (b *BigFloatArray) SetLength(length int) {
	if cap(*b) > length {
		*b = (*b)[:length]
	} else {
		*b = make([]big.Float, length)
	}
}

func (b *BigFloatArray) Resize(length int) {
	if cap(*b) > length {
		oldLength := len(*b)
		*b = (*b)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*b)[i].SetInt64(0)
			}
		}
	} else {
		tempSlice := make([]big.Float, length)
		copy(tempSlice, *b)
		*b = tempSlice
	}
}

func (b *BigFloatArray) Push(item big.Float) {
	*b = append(*b, item)
}

func (b *BigFloatArray) Pop() *big.Float {
	removed := (*b)[len(*b)-1]
	(*b)[len(*b)-1] = big.Float{}
	*b = (*b)[:len(*b)-1]
	return &removed
}

func (b *BigFloatArray) IsEmpty() bool {
	return len(*b) == 0
}

func (b *BigFloatArray) Clear() {
	*b = (*b)[:0]
}

type BigRatArray []big.Rat

const RatArrayName = packageName + ".rat"

func (b *BigRatArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(RatArrayName, "rat", "readable/writable Rat array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*b))
		} else {
			b.SetLength(arrayInspector.GetLength())
		}
		for index := range *b {
			arrayInspector.Rat(&(*b)[index], -1)
		}
		arrayInspector.End()
	}
}

func (b *BigRatArray) SetLength(length int) {
	if cap(*b) > length {
		*b = (*b)[:length]
	} else {
		*b = make([]big.Rat, length)
	}
}

func (b *BigRatArray) Resize(length int) {
	if cap(*b) > length {
		oldLength := len(*b)
		*b = (*b)[:length]
		if oldLength < length {
			for i := oldLength; i < length; i++ {
				(*b)[i].SetInt64(0)
			}
		}
	} else {
		tempSlice := make([]big.Rat, length)
		copy(tempSlice, *b)
		*b = tempSlice
	}
}

func (b *BigRatArray) Push(item big.Rat) {
	*b = append(*b, item)
}

func (b *BigRatArray) Pop() *big.Rat {
	removed := (*b)[len(*b)-1]
	(*b)[len(*b)-1] = big.Rat{}
	*b = (*b)[:len(*b)-1]
	return &removed
}

func (b *BigRatArray) IsEmpty() bool {
	return len(*b) == 0
}

func (b *BigRatArray) Clear() {
	*b = (*b)[:0]
}
