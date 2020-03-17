package common

import ()

type SimpleCallback func()
type GenericCallback func(interface{})
type ErrorCallback func(error)
type StringCallback func(string)

func (s SimpleCallback) Call() {
	if s != nil {
		s()
	}
}

func (g GenericCallback) Call(value interface{}) {
	if g != nil {
		g(value)
	}
}

func (e ErrorCallback) Call(err error) {
	if e != nil {
		e(err)
	}
}

func (s StringCallback) Call(value string) {
	if s != nil {
		s(value)
	}
}
