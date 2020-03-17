package fromjson

import (
	"bytes"
)

var nullValue = []byte("null")

func isNull(data []byte) bool {
	return bytes.Compare(data, nullValue) == 0
}
