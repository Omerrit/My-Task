package jsonrpc

type Error int

const (
	ErrUnsupportedJsonRpcVersion Error = -32769
	ErrParse                     Error = -32700
	ErrInvalidRequest            Error = -32600
	ErrMethodNotFound            Error = -32601
	ErrInvalidParams             Error = -32602
	ErrInternalError             Error = -32603
	ErrEmpty                     Error = -30000
)

var serverErrorMessages = map[Error]string{
	ErrUnsupportedJsonRpcVersion: "unsupported json rpc version",
	ErrParse:                     "parse error",
	ErrInvalidRequest:            "invalid request",
	ErrMethodNotFound:            "method not found",
	ErrInvalidParams:             "invalid params",
	ErrInternalError:             "internal error",
	ErrEmpty:                     ""}

func (j Error) Error() string {
	return serverErrorMessages[j]
}
