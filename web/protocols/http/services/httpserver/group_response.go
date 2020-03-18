package httpserver

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/jsonrpc"
)

type groupResponse []jsonrpc.GenericResponse

const groupResponseName = packageName + ".groupresp"

func (g *groupResponse) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(groupResponseName, jsonrpc.GenericResponseName, "")
	{
		if arrayInspector.IsReading() {
			return
		}
		arrayInspector.SetLength(len(*g))
		for index := range *g {
			objectInspector := arrayInspector.Value().Object(jsonrpc.GenericResponseName, "")
			(*g)[index].Embed(objectInspector)
			objectInspector.End()
		}
		arrayInspector.End()
	}
}
