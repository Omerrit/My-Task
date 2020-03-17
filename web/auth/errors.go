package auth

import (
	"gerrit-share.lan/go/basicerrors"
)

type authError int

const (
	ErrWrongIdLen authError = iota
	ErrNoConnectionId
	ErrNoUser
	ErrWrongIdType
)

var authErrorMessages = map[authError]string{
	ErrWrongIdLen:     "Wrong connection id length when trying to read it",
	ErrNoConnectionId: "Unknown connection id in request",
	ErrNoUser:         "No user was bound to this connection",
	ErrWrongIdType:    "Invalid source type for id (needs []byte), scanner error"}

var authErrorCategory = map[authError]basicerrors.BasicError{
	ErrNoUser: basicerrors.Forbidden}

func (e authError) Error() string {
	return authErrorMessages[e]
}

func (e authError) Unwrap() error {
	return authErrorCategory[e]
}
