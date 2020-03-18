package jsonrpc

import (
	"fmt"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectable"
)

const Version = "2.0"
const packageName = "jsonrpc"

type Request struct {
	JsonRpc string
	Method  string
	Params  inspect.Inspectable
}

func (r *Request) Embed(i *inspect.ObjectInspector) {
	if !i.IsReading() {
		return
	}
	i.String(&r.JsonRpc, "jsonrpc", true, "json rpc version, MUST be exactly 2.0")
	i.String(&r.Method, "method", true, "method name")
}

type RequestParams Request

func (r *RequestParams) Embed(i *inspect.ObjectInspector) {
	if !i.IsReading() {
		return
	}
	r.Params.Inspect(i.Value("params", false, "a structured value that holds the parameter values to be used during the invocation of the method"))
}

type Response struct {
	JsonRpc string
	Result  GenericResponse
}

const ResponseName = packageName + ".resp"

func NewResponse(typeId inspect.TypeId) *Response {
	return &Response{
		JsonRpc: Version,
		Result: GenericResponse{
			Result: inspectable.NewGenericValue(typeId)}}
}

func (r *Response) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(ResponseName, "json rpc response object")
	{
		if objectInspector.IsReading() {
			return
		}
		objectInspector.String(&r.JsonRpc, "jsonrpc", true, "json rpc version, MUST be exactly 2.0")
		r.Result.Embed(objectInspector)
		objectInspector.End()
	}
}

type jsonRpcError struct {
	code    int
	message string
	err     error
}

const jsonRpcErrorName = packageName + ".error"

func (e *jsonRpcError) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(jsonRpcErrorName, "json rpc error object")
	{
		objectInspector.Int(&e.code, "code", true, "error code")
		objectInspector.String(&e.message, "message", true, "error message")
		objectInspector.End()
	}
}

func (e *jsonRpcError) Error() string {
	return fmt.Sprintf("%v: %v", e.code, e.message)
}

func (e *jsonRpcError) Unwrap() error {
	return e.err
}

func (e *jsonRpcError) Is(err error) bool {
	code, ok := err.(Error)
	return ok && int(code) == e.code
}

func (e *jsonRpcError) As(target interface{}) bool {
	jrpcErr, ok := target.(*Error)
	if !ok {
		return false
	}
	*jrpcErr = Error(e.code)
	return true
}

func Describe(err error, code Error) *jsonRpcError {
	e := &jsonRpcError{
		code: int(code),
		message: func() string {
			if err == nil {
				return code.Error()
			} else {
				return fmt.Sprintf("%v: %v", code.Error(), err.Error())
			}
		}(),
		err: err,
	}
	return e
}

func ErrorFromCode(code Error) *jsonRpcError {
	return Describe(nil, code)
}

type ResponseBatch []*Response

const ResponseBatchName = packageName + ".respbatch"

func (b *ResponseBatch) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(ResponseBatchName, ResponseName, "json rpc response batch")
	{
		for index := range *b {
			(*b)[index].Inspect(arrayInspector.Value())
		}
		arrayInspector.End()
	}
}
