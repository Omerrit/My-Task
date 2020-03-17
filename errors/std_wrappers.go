package errors

import (
	"errors"
	"fmt"
)

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Errorf(format string, args ...interface{}) StackTraceError {
	err := fmt.Errorf(format, args...)
	unwrapped := errors.Unwrap(err)
	if unwrapped != nil {
		if stackTraceErr, ok := unwrapped.(StackTraceError); ok {
			return &wrappedStringStackTraceError{stackTraceErr, err.Error()}
		}
		return &wrappedStringError{makeStringError(err.Error(), 1), unwrapped}
	}
	strerr := makeStringError(err.Error(), 1)
	return &strerr
}
