package treeedit

import (
	"gerrit-share.lan/go/basicerrors"
)

type myErr int

const (
	ErrInvalidParent myErr = iota
	ErrInvalidPosition
	ErrInvalidDivision
	ErrInvalidId
	ErrCannotRemoveRoot
	ErrPositionIsRoot
	ErrNothingToSave
)

var myErrMessages = map[myErr]string{
	ErrInvalidParent:    "Invalid parent id",
	ErrInvalidPosition:  "Invalid position id",
	ErrInvalidDivision:  "Invalid division id",
	ErrInvalidId:        "Invalid id",
	ErrCannotRemoveRoot: "Root item can't be removed",
	ErrPositionIsRoot:   "Position is root, something unimaginable happened",
	ErrNothingToSave:    "Nothing to save"}

var myErrWrapped = map[myErr]error{
	ErrInvalidParent:    basicerrors.NotFound,
	ErrInvalidPosition:  basicerrors.NotFound,
	ErrInvalidDivision:  basicerrors.NotFound,
	ErrInvalidId:        basicerrors.NotFound,
	ErrCannotRemoveRoot: basicerrors.Forbidden}

func (m myErr) Error() string {
	return myErrMessages[m]
}

func (m myErr) Unwrap() error {
	return myErrWrapped[m]
}
