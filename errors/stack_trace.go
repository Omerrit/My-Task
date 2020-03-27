package errors

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	defaultPCSize = 16
	maxPCSize     = 64 * defaultPCSize
)

type StackFrame struct {
	File     string
	Line     int
	Function string
	Package  string
	Offset   int
}

func (sf StackFrame) String() string {
	return fmt.Sprintf("%v %v\n\t%v:%v (+0x%x)", sf.Package, sf.Function, sf.File, sf.Line, sf.Offset)
}

type CallStack []StackFrame

func (cs *CallStack) CutBottom(isCutOff func(StackFrame) bool) {
	for i := len(*cs) - 1; i > -1; i-- {
		if !isCutOff((*cs)[i]) {
			*cs = (*cs)[:i+1]
			break
		}
	}
}

func (cs *CallStack) CutTop(isCutOff func(StackFrame) bool) {
	for i, val := range *cs {
		if !isCutOff(val) {
			*cs = (*cs)[i:]
			break
		}
	}
}

func (cs CallStack) String() string {
	var builder strings.Builder
	for _, val := range cs {
		builder.WriteString(val.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

func splitFunction(function string) (string, string) {
	index := strings.LastIndex(function, "/") + 1
	indexOfDot := strings.Index(function[index:], ".")
	if indexOfDot < 0 {
		return function, ""
	}
	return function[:index+indexOfDot], function[index+indexOfDot+1:]
}

func CallerStack(depth uint) CallStack {
	if depth == 0 {
		return nil
	}
	return getStackTrace(2, depth)
}

func MyStack(depth uint) CallStack {
	if depth == 0 {
		return nil
	}
	return getStackTrace(1, depth)
}

func getStackTrace(skip uint, depth uint) CallStack {
	pcSize := defaultPCSize
	var pc []uintptr
	var n int
	for {
		pc = make([]uintptr, pcSize)
		n = runtime.Callers(int(2+skip), pc)
		if n >= int(depth) {
			n = int(depth + 1)
			break
		}
		if n != len(pc) || n == maxPCSize {
			break
		}
		pcSize = pcSize * 2
	}

	frames := runtime.CallersFrames(pc[:n])

	var tempFrames CallStack
	for {
		frame, more := frames.Next()
		packageName, function := splitFunction(frame.Function)
		tempFrames = append(tempFrames, StackFrame{File: frame.File, Line: frame.Line, Function: function, Package: packageName, Offset: int(frame.PC - frame.Entry + 1)})
		if !more {
			break
		}
	}
	return tempFrames
}
