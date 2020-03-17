package flags

import (
	"strconv"
)

// -- int32 Value
type int32Value int32

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 10, 32)
	*i = int32Value(v)
	return err
}

func (i *int32Value) String() string { return strconv.Itoa(int(*i)) }
