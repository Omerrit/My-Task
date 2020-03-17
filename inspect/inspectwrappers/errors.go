package inspectwrappers

import ()

type wrapperErr int

const (
	ErrCantWrite wrapperErr = iota
	ErrAmbiguity
	ErrUnsupported
	ErrWrongType
)

var wrapperErrNames = map[wrapperErr]string{
	ErrCantWrite:   "Trying to write into read only value",
	ErrAmbiguity:   "byte slice type is not supported in auto guessing due to ambiguity",
	ErrUnsupported: "unsupported type",
	ErrWrongType:   "new value should be of the same type as the one being replaced"}

func (w wrapperErr) Error() string {
	return wrapperErrNames[w]
}
