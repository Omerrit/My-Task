package simpleutils

import ()

func ResizeBytes(bytes []byte, size int) []byte {
	if cap(bytes) < size {
		return make([]byte, len(bytes), size)
	}
	return bytes[:size]
}
