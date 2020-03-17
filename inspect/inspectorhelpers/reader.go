package inspectorhelpers

import (
	"gerrit-share.lan/go/inspect"
	"math/big"
)

type Reader struct {
}

func (*Reader) IsReading() bool {
	return true
}

func (*Reader) MapWriteInt(key string, value int) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteInt32(key string, value int32) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteInt64(key string, value int64) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteFloat32(key string, value float32, format byte, precision int) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteFloat64(key string, value float64, format byte, precision int) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteString(key string, value string) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteByteString(key string, value []byte) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteBytes(key string, value []byte) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteBool(key string, value bool) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteBigInt(key string, value *big.Int) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteRat(key string, value *big.Rat, precision int) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteBigFloat(key string, value *big.Float, format byte, precision int) error {
	return inspect.ErrWritingToReader
}
func (*Reader) MapWriteValue(key string) error {
	return inspect.ErrWritingToReader
}
