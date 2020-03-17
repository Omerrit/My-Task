package fromjson

import (
	"gerrit-share.lan/go/basicerrors"
	"strconv"
)

type parseError int

const (
	ErrShouldBeNull parseError = iota
	ErrNotBool
	ErrWrongLength
	ErrPropertyNotFound
	ErrParsingValue
	ErrStringRequired
	ErrIntegerRequired
	ErrFloatRequired
	ErrBoolRequired
	ErrBigIntRequired
	ErrRatRequired
	ErrBigFloatRequired
	ErrObjectRequired
	ErrArrayRequired
)

var parseErrorMessages = map[parseError]string{
	ErrShouldBeNull:     "parameter is not null",
	ErrNotBool:          "boolean parameter expected",
	ErrWrongLength:      "array reader was requested more elements than it can provide",
	ErrPropertyNotFound: "mandatory object property is abcent",
	ErrParsingValue:     "failed to parse value property",
	ErrStringRequired:   "not a string",
	ErrFloatRequired:    "not a floating point value",
	ErrIntegerRequired:  "not an integer value",
	ErrBoolRequired:     "not a boolean value",
	ErrBigIntRequired:   "not a large integer value (should be passed as a string)",
	ErrRatRequired:      "not a large fixed point value (should be passed as a string)",
	ErrBigFloatRequired: "not a large floating point value (should be passed as a string)",
	ErrObjectRequired:   "not a json object",
	ErrArrayRequired:    "not an array"}

func (p parseError) Error() string {
	return parseErrorMessages[p]
}

func (p parseError) Unwrap() error {
	return basicerrors.BadParameter
}

type fullParseError struct {
	err    error
	offset int
}

func (p *fullParseError) Error() string {
	return "parsing error at offset " + strconv.Itoa(p.offset) + ": " + p.err.Error()
}

func (p *fullParseError) Unwrap() error {
	return p.err
}

func makeParseError(offset int, err error) error {
	if pe, ok := err.(*fullParseError); ok {
		pe.offset += offset
		return pe
	}
	return &fullParseError{err, offset}
}
