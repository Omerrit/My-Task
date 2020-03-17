package errors

import ()

type StackTraceError interface {
	error
	Unwrap() error
	StackTrace() CallStack
	replaceStackTrace(CallStack)
}

type DescribedError interface {
	StackTraceError
	Descriptions() []string
	appendDescription(string)
}
