package errors

import (
	"fmt"
	"strings"
)

func makeStringError(message string, skipTopStackItems uint) stringError {
	stack := getStackTrace(skipTopStackItems+1, 0)
	if len(stack) > 0 && (strings.Contains(stack[0].Function, "init.") || stack[0].Function == "init") {
		return stringError{body: message}
	}
	return stringError{body: message, stackTrace: stackTrace{stack}}
}

func New(errstr string) StackTraceError {
	err := makeStringError(errstr, 1)
	return &err
}

func FromPrintln(values ...interface{}) StackTraceError {
	err := makeStringError(fmt.Sprintln(values...), 1)
	return &err
}

func describe(err error, description string, skipTopStackItems uint) DescribedError {
	if err == nil {
		return nil
	}
	described, ok := err.(DescribedError)
	if !ok {
		described = &describedError{StackTraceError: wrap(err, skipTopStackItems+1)}
	}
	described.appendDescription(description)
	return described
}

func Describe(err error, description string) DescribedError {
	return describe(err, description, 1)
}

func Describef(err error, format string, args ...interface{}) DescribedError {
	return describe(err, fmt.Sprintf(format, args...), 1)
}

func simplifyArray(err error) error {
	array, ok := err.(ErrorArray)
	if !ok {
		return err
	}
	return array.ToError()
}

func wrapAndReplaceStack(err error, newStack CallStack) StackTraceError {
	err = simplifyArray(err)
	ste, ok := err.(StackTraceError)
	if ok {
		ste.replaceStackTrace(newStack)
		return ste
	}
	return &wrappedError{stackTrace: stackTrace{newStack}, error: err}
}

func wrap(err error, skipTopStackItems uint) StackTraceError {
	err = simplifyArray(err)
	if err == nil {
		return nil
	}
	stackTraceErr, ok := err.(StackTraceError)
	if ok {
		return stackTraceErr
	}
	return &wrappedError{stackTrace: stackTrace{getStackTrace(skipTopStackItems+1, 0)}, error: err}
}

func Wrap(err error) StackTraceError {
	return wrap(err, 1)
}
