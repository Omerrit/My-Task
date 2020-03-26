package replies

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/cookies"
	"time"
)

const packageName = "httpreplies"

type HttpResponse struct {
	Response  []byte
	SessionId []byte
}

const HttpResponseName = packageName + ".httpresp"

func NewHttpResponse(info cookies.CookieInfo) *HttpResponse {
	response := &HttpResponse{}
	if float64((info.ExpiresAt-time.Now().Unix())*1e9)/float64(info.SessionDuration) > info.SessionReset {
		return response
	}
	sessionId, _ := info.SessionId.MarshalText()
	response.SessionId = sessionId
	return response
}

func (h *HttpResponse) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(HttpResponseName, "")
	{
		objectInspector.ByteString(&h.Response, "response", true, "response body")
		objectInspector.ByteString(&h.SessionId, "sessionid", false, "session id")
		objectInspector.End()
	}
}

func (h HttpResponse) Visit(visitor actors.ResponseVisitor) {
	visitor.Reply(&h)
}

func init() {
	inspectables.Register(HttpResponseName, func() inspect.Inspectable { return new(HttpResponse) })
}

type HttpResponsePromise struct {
	actors.ResponsePromise
}

func (h *HttpResponsePromise) Deliver(value *HttpResponse) {
	h.DeliverUntyped(value)
}
