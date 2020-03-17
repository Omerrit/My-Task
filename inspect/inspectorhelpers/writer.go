package inspectorhelpers

import (
	"gerrit-share.lan/go/inspect"
	"math/big"
)

type Writer struct {
}

func (*Writer) IsReading() bool {
	return false
}

func (*Writer) MapNextKey() (string, error) {
	return "", inspect.ErrReadingFromWriter
}
func (*Writer) MapReadInt() (int, error) {
	return 0, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadInt32() (int32, error) {
	return 0, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadInt64() (int64, error) {
	return 0, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadFloat32() (float32, error) {
	return 0, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadFloat64() (float64, error) {
	return 0, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadString() (string, error) {
	return "", inspect.ErrReadingFromWriter
}
func (*Writer) MapReadByteString(value []byte) ([]byte, error) {
	return value, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadBytes(value []byte) ([]byte, error) {
	return value, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadBool() (bool, error) {
	return false, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadBigInt(value *big.Int) (*big.Int, error) {
	return value, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadRat(value *big.Rat) (*big.Rat, error) {
	return value, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadBigFloat(value *big.Float) (*big.Float, error) {
	return value, inspect.ErrReadingFromWriter
}
func (*Writer) MapReadValue() error {
	return inspect.ErrReadingFromWriter
}
