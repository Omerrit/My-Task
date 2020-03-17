package frombytes

import "gerrit-share.lan/go/basicerrors"

type fromJsonError int

const (
	ErrDuplicatePath fromJsonError = iota
	ErrTooFewParameters
)

var serverErrorMessages = map[fromJsonError]string{
	ErrDuplicatePath:    "\"path\" parameter has been detected both in url and json",
	ErrTooFewParameters: "too few parameters"}

func (m fromJsonError) Error() string {
	return serverErrorMessages[m]
}

func (m fromJsonError) Unwrap() error {
	return basicerrors.BadParameter
}
