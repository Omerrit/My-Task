package jsonrpc

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectable"
)

type GenericResponse struct {
	Result *inspectable.GenericValue
	Err    error
}

func (g *GenericResponse) Inspect(objectInspector *inspect.ObjectInspector) {
	if objectInspector.IsReading() {
		return
	}
	if g.Err != nil {
		if e, ok := g.Err.(inspect.Inspectable); ok {
			e.Inspect(objectInspector.Value("error", false, "error"))
			return
		}
		s := g.Err.Error()
		objectInspector.String(&s, "error", false, "error")
		return
	}
	g.Result.Inspect(objectInspector.Value("result", false, "result"))
}
