package treeedit

import (
	"encoding/base64"
	"encoding/binary"
)

func EncodeIntId(id int64) string {
	var buffer [binary.MaxVarintLen64]byte
	return base64.RawURLEncoding.EncodeToString(buffer[:binary.PutVarint(buffer[:], id)])
}
