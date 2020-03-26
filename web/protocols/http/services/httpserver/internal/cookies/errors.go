package cookies

import "gerrit-share.lan/go/basicerrors"

type tokenError int

const (
	ErrIncorrectSessionToken tokenError = iota
	ErrTokenExpired
)

var serverErrorMessages = map[tokenError]string{
	ErrIncorrectSessionToken: "incorrect session token",
	ErrTokenExpired:          "token expired"}

var serverErrorWrapped = map[tokenError]error{
	ErrIncorrectSessionToken: basicerrors.Forbidden,
	ErrTokenExpired:          basicerrors.Forbidden}

func (m tokenError) Error() string {
	return serverErrorMessages[m]
}

func (m tokenError) Unwrap() error {
	return serverErrorWrapped[m]
}
