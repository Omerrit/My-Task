package httpserver

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/jsonrpc"
)

type groupResponse []jsonrpc.GenericResponse

func (g *groupResponse) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array("", "", "")
	if arrayInspector.IsReading() {
		return
	}
	arrayInspector.SetLength(len(*g))
	for index := range *g {
		objectInspector := arrayInspector.Value().Object("", "")
		(*g)[index].Inspect(objectInspector)
		objectInspector.End()
	}
	arrayInspector.End()
}
