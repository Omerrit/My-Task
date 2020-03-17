package httpserver

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"net/http"
)

const (
	urlPathName            = packageName + ".doc"
	httpRequestName        = packageName + ".request"
	editEndpointName       = packageName + ".editendpoint"
	editEndpointByNameName = packageName + ".editendpointbyname"
	getEndpointsName       = packageName + ".getendpoints"
)

type urlPath struct {
	id string
}

func (o *urlPath) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(urlPathName, "get api documentation")
	objectInspector.String(&o.id, "path", false, "resource path")
	objectInspector.End()
}

type httpRequest http.Request

func (h *httpRequest) Inspect(i *inspect.GenericInspector) {}

type editEndpoint struct {
	resource      string
	method        string
	httpMethod    string
	oldResource   string
	oldMethod     string
	oldHttpMethod string
}

func (e *editEndpoint) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(editEndpointName, "edit endpoint")
	{
		objectInspector.String(&e.resource, "resource", true, "new resource name")
		objectInspector.String(&e.method, "method", true, "new method name")
		objectInspector.String(&e.httpMethod, "http", true, "new http method name")
		objectInspector.String(&e.oldResource, "old_resource", false, "old resource name")
		objectInspector.String(&e.oldMethod, "old_method", false, "old method name")
		objectInspector.String(&e.oldHttpMethod, "old_http", false, "old http method name")
	}
	objectInspector.End()
}

type editEndpointByName struct {
	resource   string
	method     string
	httpMethod string
	name       string
}

func (e *editEndpointByName) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(editEndpointByNameName, "edit endpoint by original name")
	{
		objectInspector.String(&e.resource, "resource", true, "new resource name")
		objectInspector.String(&e.method, "method", true, "new method name")
		objectInspector.String(&e.httpMethod, "http", true, "new http method name")
		objectInspector.String(&e.name, "name", false, "original name")
	}
	objectInspector.End()
}

type getEndpoints struct {
	edited bool
}

func (g *getEndpoints) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(getEndpointsName, "get endpoints")
	{
		objectInspector.Bool(&g.edited, "edited", false, "only edited / not edited endpoints")
	}
	objectInspector.End()
}

func init() {
	inspectables.Register(urlPathName, func() inspect.Inspectable { return new(urlPath) })
	inspectables.Register(httpRequestName, func() inspect.Inspectable { return new(httpRequest) })
	inspectables.Register(editEndpointName, func() inspect.Inspectable { return new(editEndpoint) })
	inspectables.Register(editEndpointByNameName, func() inspect.Inspectable { return new(editEndpointByName) })
	inspectables.Register(getEndpointsName, func() inspect.Inspectable { return new(getEndpoints) })
}
