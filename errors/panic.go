package errors

import (
	"encoding/json"
	"fmt"
)

type UserPanic struct {
	StackTraceError
}

type GoPanic struct {
	StackTraceError
}

func isNotPanic(frame StackFrame) bool {
	if frame.Package == "runtime" && frame.Function == "gopanic" {
		return false
	}
	return true
}

func isPanic(frame StackFrame) bool {
	if frame.Package == "runtime" && frame.Function == "gopanic" {
		return true
	}
	return false
}

func isRuntime(frame StackFrame) bool {
	if frame.Package == "runtime" {
		return true
	}
	return false
}

func generatePanicErr(src error) StackTraceError {
	stack := MyStack(maxPCSize)
	stack.CutTop(isNotPanic)
	lenBefore := len(stack)
	stack.CutTop(isRuntime)
	err := wrapAndReplaceStack(src, stack)
	if len(stack) < (lenBefore - 1) {
		return &GoPanic{err}
	}
	return &UserPanic{err}
}

func dataToString(data interface{}) string {
	result, err := json.Marshal(data)
	if err == nil {
		return string(result)
	} else {
		return fmt.Sprintf("%#v", data)
	}
}

func RecoverToError(src interface{}) StackTraceError {
	if src == nil {
		return nil
	}

	var newErr StackTraceError
	switch val := src.(type) {
	case string:
		newErr = &stringError{body: val}
	case error:
		newErr = &wrappedError{error: val}
	default:
		newErr = &unknownError{
			stringError: stringError{body: dataToString(val)},
			body:        val}
	}

	return generatePanicErr(newErr)
}
