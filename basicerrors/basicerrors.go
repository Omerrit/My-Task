package basicerrors

import ()

type BasicError int

const (
	NotFound BasicError = iota
	BadParameter
	Forbidden
	ProcessorDead
)

var basicErrors = map[BasicError]string{
	NotFound:      "not found",
	BadParameter:  "bad prameter",
	Forbidden:     "forbidden",
	ProcessorDead: "dead"}

func (b BasicError) Error() string {
	return basicErrors[b]
}
