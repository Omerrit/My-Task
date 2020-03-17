package arrays

import (
	"gerrit-share.lan/go/inspect"
	"math/big"
)

const (
	BigIntArrayName   = packageName + ".bigint"
	RatArrayName      = packageName + ".rat"
	BigFloatArrayName = packageName + ".bigfloat"
)

type BigIntArray []*big.Int

func (s *BigIntArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(BigIntArrayName, "bigInt", "readable/writable bigInt array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.BigInt((*s)[index])
	}
	arrayInspector.End()
}

func (s *BigIntArray) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([]*big.Int, length)
		for i := 0; i < len(*s); i++ {
			(*s)[i] = big.NewInt(0)
		}
	}
}

type BigFloatArray []*big.Float

func (s *BigFloatArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(BigFloatArrayName, "bigFloat", "readable/writable bigFloat array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.BigFloat((*s)[index], 'g', -1)
	}
	arrayInspector.End()
}

func (s *BigFloatArray) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([]*big.Float, length)
		for i := 0; i < len(*s); i++ {
			(*s)[i] = big.NewFloat(0)
		}
	}
}

type BigRatArray []*big.Rat

func (s *BigRatArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(RatArrayName, "rat", "readable/writable Rat array")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*s))
	} else {
		s.init(arrayInspector.GetLength())
	}
	for index := range *s {
		arrayInspector.Rat((*s)[index], -1)
	}
	arrayInspector.End()
}

func (s *BigRatArray) init(length int) {
	if cap(*s) > length {
		*s = (*s)[:length]
	} else {
		*s = make([]*big.Rat, length)
		for i := 0; i < len(*s); i++ {
			(*s)[i] = big.NewRat(0, 1)
		}
	}
}
