package httpserver

import (
	"errors"
	"fmt"
	"gerrit-share.lan/go/basicerrors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/json/tojson"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/jsonrpc"
	"net/http"
)

var errorToHttpCodeMap = map[basicerrors.BasicError]int{
	basicerrors.BadParameter:  http.StatusBadRequest,
	basicerrors.NotFound:      http.StatusNotFound,
	basicerrors.Forbidden:     http.StatusForbidden,
	basicerrors.ProcessorDead: http.StatusServiceUnavailable}

func errorToHttpCode(err error) int {
	var category basicerrors.BasicError
	if !errors.As(err, &category) {
		return http.StatusInternalServerError
	}
	return errorToHttpCodeMap[category]
}

var rpcErrorToHttpCodeMap = map[jsonrpc.Error]int{
	jsonrpc.ErrMethodNotFound:            http.StatusMethodNotAllowed,
	jsonrpc.ErrParse:                     http.StatusBadRequest,
	jsonrpc.ErrInvalidRequest:            http.StatusBadRequest,
	jsonrpc.ErrInvalidParams:             http.StatusBadRequest,
	jsonrpc.ErrUnsupportedJsonRpcVersion: http.StatusBadRequest,
	jsonrpc.ErrInternalError:             http.StatusInternalServerError}

func rpcErrorToHttpCode(err error) int {
	var category jsonrpc.Error
	if !errors.As(err, &category) {
		return http.StatusInternalServerError
	}
	return rpcErrorToHttpCodeMap[category]
}

func processError(writer http.ResponseWriter, err error) {
	writeResponseStatus(writer, err)
	writeError(writer, err)
}

func writeResponseStatus(writer http.ResponseWriter, err error) {
	if errors.As(err, new(jsonrpc.Error)) {
		writer.WriteHeader(rpcErrorToHttpCode(err))
	} else {
		if err == ErrMethodNotAllowed {
			writer.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			writer.WriteHeader(errorToHttpCode(err))
		}
	}
}

func writeError(writer http.ResponseWriter, err error) {
	if inspectableErr, ok := err.(inspect.Inspectable); ok {
		i := &tojson.Inspector{}
		serializer := inspect.NewGenericInspector(i)
		inspectableErr.Inspect(serializer)
		writer.Write(i.Output())
		return
	}
	writer.Write([]byte(fmt.Sprintf("%#v, %v", err, err)))
}
