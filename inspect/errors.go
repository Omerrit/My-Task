package inspect

import "gerrit-share.lan/go/basicerrors"

type serializerError int

const (
	ErrMandatoryFieldAbsent serializerError = iota
	ErrWritingToReader
	ErrReadingFromWriter
)

var serializerErrorMessages = map[serializerError]string{
	ErrMandatoryFieldAbsent: "mandatory field is absent",
	ErrWritingToReader:      "writing to a reading inspector",
	ErrReadingFromWriter:    "reading from a writing inspector"}

var serializerErrorCategory = map[serializerError]basicerrors.BasicError{
	ErrMandatoryFieldAbsent: basicerrors.BadParameter}

func (e serializerError) Error() string {
	return serializerErrorMessages[e]
}

func (e serializerError) Unwrap() error {
	err, ok := serializerErrorCategory[e]
	if ok {
		return err
	} else {
		return nil
	}
}
