package errors

import (
	"testing"
)

var testErr = New("hello")
var testErr2 StackTraceError
var testErrFold StackTraceError
var testMadeErr = makeOtherErr()

func makeErr() {
	testErrFold = New("nope")
}

func init() {
	makeErr()
	testErr2 = New("hi")
}

func makeOtherErr() StackTraceError {
	return New("hi")
}

func TestGlobalNew(t *testing.T) {
	if len(testErr.StackTrace()) > 0 {
		t.Error("stack should be empty in error initialized in global scope", testErr, testErr.StackTrace())
	}
	t.Log(testErr.StackTrace())
}

func TestInitNew(t *testing.T) {
	if len(testErr2.StackTrace()) > 0 {
		t.Error("stack should be empty in error initialized inside init", testErr2, testErr2.StackTrace())
	}
	t.Log(testErr2.StackTrace())
}

func TestInitFoldNew(t *testing.T) {
	if len(testErrFold.StackTrace()) == 0 {
		t.Error("stack should't be empty in errors initialized inside a function that init() calls", testErrFold, testErrFold.StackTrace())
	}
	t.Log(testErrFold.StackTrace())
}

func TestMadeError(t *testing.T) {
	if len(testMadeErr.StackTrace()) == 0 {
		t.Error("stack should't be empty in errors initialized inside a function that init() calls", testMadeErr, testMadeErr.StackTrace())
	}
	t.Log(testMadeErr.StackTrace())
}
