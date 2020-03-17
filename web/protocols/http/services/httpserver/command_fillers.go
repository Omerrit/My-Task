package httpserver

import (
	"gerrit-share.lan/go/basicerrors"
	"gerrit-share.lan/go/debug"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"gerrit-share.lan/go/web/protocols/http/serializers"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/frombytes"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/metadata"
	"net/http"
	"net/url"
)

const maxMemoryForMultipartForm = 1024

func commandFromUrl(values url.Values, generator inspectables.Creator) (inspect.Inspectable, error) {
	command := generator()
	deserializer := &serializers.FromUrl{Values: values}
	i := inspect.NewGenericInspector(deserializer)
	command.Inspect(i)
	if i.GetError() != nil {
		return nil, basicerrors.Augment(i.GetError(), basicerrors.BadParameter)
	}
	return command, nil
}

func commandFromQueryString(request *http.Request,
	endpoint endpointInfo, path []byte) (inspect.Inspectable, error) {
	err := request.ParseForm()
	if err != nil {
		return nil, basicerrors.Augment(err, basicerrors.BadParameter)
	}
	debug.Printf("calling %s, path=%s\n", endpoint.OriginalName, string(path))
	err = modifyFormIfSingleKey(request.Form, endpoint)
	if err != nil {
		return nil, err
	}
	if !(endpoint.CommandMetaData.PathIndex == metadata.NoPath) && len(path) > 0 {
		if len(request.Form.Get("path")) > 0 {
			return nil, frombytes.ErrDuplicatePath
		}
		request.Form.Add("path", string(path))
	}
	return commandFromUrl(request.Form, endpoint.CommandGenerator)
}

func modifyFormIfSingleKey(form url.Values, endpoint endpointInfo) error {
	if len(form) == 1 {
		params := endpoint.CommandMetaData.UnderlyingValues
		var notPathIndex int
		var value string
		for key := range form {
			value = key
		}
		if len(params) > 2 || ((endpoint.CommandMetaData.PathIndex == metadata.NoPath) && len(params) > 1) {
			return frombytes.ErrTooFewParameters
		}
		if len(params) == 2 && !(endpoint.CommandMetaData.PathIndex == metadata.NoPath) {
			notPathIndex = (endpoint.CommandMetaData.PathIndex + 1) % len(params)
		}
		if len(form[value]) == 1 && form[value][0] == "" && params[0].Name != value {
			form.Add(params[notPathIndex].Name, value)
			form.Del(value)
		}
	}
	return nil
}

func commandFromFormData(request *http.Request,
	endpoint endpointInfo, path []byte) (inspect.Inspectable, error) {
	err := request.ParseMultipartForm(maxMemoryForMultipartForm)
	if err != nil {
		return nil, basicerrors.Augment(err, basicerrors.BadParameter)
	}
	if !(endpoint.CommandMetaData.PathIndex == metadata.NoPath) && len(path) > 0 {
		if len(request.MultipartForm.Value["path"]) > 0 {
			return nil, frombytes.ErrDuplicatePath
		}
		request.MultipartForm.Value["path"] = append(request.MultipartForm.Value["path"], string(path))
	}
	return commandFromUrl(request.MultipartForm.Value, endpoint.CommandGenerator)
}

func commandFromPath(endpoint endpointInfo, path []byte) (inspect.Inspectable, error) {
	values := make(url.Values)
	if !(endpoint.CommandMetaData.PathIndex == metadata.NoPath) {
		values.Add("path", string(path))
	}
	return commandFromUrl(values, endpoint.CommandGenerator)
}
