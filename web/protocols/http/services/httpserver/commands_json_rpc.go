package httpserver

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/jsonrpc"
)

type requestWithDestination struct {
	request      jsonrpc.Request
	destination  actors.ActorService
	responseType inspect.TypeId
	err          error
}

const requestWithDestinationName = packageName + ".rpcreq"

type rpcRequestBatch struct {
	data      []requestWithDestination
	endpoints *HttpRestEndpoints
}

const rpcRequestBatchName = packageName + ".rpcreqbatch"

func (r *rpcRequestBatch) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(rpcRequestBatchName, requestWithDestinationName, "rpc request batch")
	{
		r.data = make([]requestWithDestination, arrayInspector.GetLength())
		for index := range r.data {
			objectInspector := arrayInspector.Value().Object(requestWithDestinationName, "")
			{
				r.data[index].request.Embed(objectInspector)
				if r.data[index].request.JsonRpc != jsonrpc.Version {
					r.data[index].err = jsonrpc.ErrorFromCode(jsonrpc.ErrUnsupportedJsonRpcVersion)
					continue
				}

				endpoint, err := r.endpoints.getEndpointByOriginalName(r.data[index].request.Method)
				if err != nil {
					r.data[index].err = jsonrpc.Describe(err, jsonrpc.ErrMethodNotFound)
					continue
				}
				r.data[index].destination = endpoint.Dest
				r.data[index].responseType = endpoint.ResultInfo.TypeId
				r.data[index].request.Params = endpoint.CommandGenerator()
				(*jsonrpc.RequestParams)(&r.data[index].request).Embed(objectInspector)
				objectInspector.End()
			}
		}
		arrayInspector.End()

	}
}
