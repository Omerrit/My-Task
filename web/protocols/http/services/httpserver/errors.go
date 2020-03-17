package httpserver

import "gerrit-share.lan/go/basicerrors"

type serverError int

const (
	ErrResourceNotFound serverError = iota
	ErrMethodNotAllowed
	ErrHttpMethodNotAllowed
	ErrPathNotAllowed
	ErrUnsupportedContentType
	ErrNoContentType
)

var serverErrorMessages = map[serverError]string{
	ErrResourceNotFound:       "resource not found",
	ErrMethodNotAllowed:       "method not allowed",
	ErrHttpMethodNotAllowed:   "http method not allowed",
	ErrUnsupportedContentType: "unsupported content type",
	ErrPathNotAllowed:         "path not allowed for provided method",
	ErrNoContentType:          "no content type provided"}

var serverErrorWrapped = map[serverError]error{
	ErrResourceNotFound:       basicerrors.NotFound,
	ErrUnsupportedContentType: basicerrors.BadParameter,
	ErrMethodNotAllowed:       basicerrors.Forbidden,
	ErrHttpMethodNotAllowed:   basicerrors.Forbidden,
	ErrPathNotAllowed:         basicerrors.BadParameter,
	ErrNoContentType:          basicerrors.BadParameter}

func (m serverError) Error() string {
	return serverErrorMessages[m]
}

func (m serverError) Unwrap() error {
	return serverErrorWrapped[m]
}

type conversionError int

const (
	ErrWrongType conversionError = iota
	ErrUnsupportedResultType
	ErrUnsupportedCommandType
)

var conversionErrorMessages = map[conversionError]string{
	ErrWrongType:              "resource not found",
	ErrUnsupportedResultType:  "unsupported type of result sample",
	ErrUnsupportedCommandType: "unsupported type of command"}

func (c conversionError) Error() string {
	return conversionErrorMessages[c]
}
