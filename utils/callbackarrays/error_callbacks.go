package callbackarrays

import (
	"gerrit-share.lan/go/common"
)

type ErrorCallbacks []common.ErrorCallback

func (e *ErrorCallbacks) Push(callback common.ErrorCallback) {
	*e = append(*e, callback)
}

func (e *ErrorCallbacks) PushNonNull(callback common.ErrorCallback) {
	if callback != nil {
		*e = append(*e, callback)
	}
}

func (e ErrorCallbacks) Run(err error) {
	for i := len(e) - 1; i >= 0; i-- {
		e[i].Call(err)
	}
}

func (e *ErrorCallbacks) Clear() {
	*e = (*e)[:0]
}
