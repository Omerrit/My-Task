package actors

import (
	"gerrit-share.lan/go/inspect"
)

type streamId int

func (id streamId) IsValid() bool {
	return id != 0
}

func (id *streamId) Increment() {
	(*id)++
}

func (id *streamId) Invalidate() {
	*id = 0
}

func (id *streamId) Inspect(inspector *inspect.GenericInspector) {
	inspector.Int((*int)(id), packageName+".streamId", "incoming stream id within a receiving actor")
}

type RequestStream interface {
	inspect.Inspectable
	setStreamRequest(streamRequest)
	getStreamRequest() streamRequest
}

type RequestStreamBase streamRequest

const RequestStreamBaseName = packageName + ".request_stream"

func (r *RequestStreamBase) Inspect(inspector *inspect.GenericInspector) {
	(*streamRequest)(r).Inspect(inspector)
}

func (r *RequestStreamBase) Fill() {}

func (r *RequestStreamBase) getStreamRequest() streamRequest {
	return *(*streamRequest)(r)
}

func (r *RequestStreamBase) setStreamRequest(req streamRequest) {
	*r = RequestStreamBase(req)
}

func (r *RequestStreamBase) Visit(v ResponseVisitor) {
	v.Reply((*streamRequest)(r))
}

type outputId struct {
	streamId
	destination ActorService
}

func (o *outputId) Inspect(inspector *inspect.GenericInspector) {
	i := inspector.Object(packageName+".outputId", "stream output address (destination+id)")
	{
		o.streamId.Inspect(i.Value("id", true, ""))
		InspectActorService(&o.destination, i.Value("destination", true, ""))
		i.End()
	}
}

type sourceId struct {
	streamId
	source ActorService
}

type streamCanSend sourceId

type streamRequest struct { //reply+request
	id     outputId
	data   inspect.Inspectable //if data is an array, len(data) is the minumum size
	maxLen int                 //automatically grow array up to this size
}

func (s *streamRequest) Inspect(inspector *inspect.GenericInspector) {
	i := inspector.Object(packageName+".streamRequest", "stream 'request data' message")
	s.id.Inspect(i.Value("id", true, ""))
	s.data.Inspect(i.Value("data", true, ""))
	i.Int(&s.maxLen, "maxLen", true, "maximum data length to fill by stream source")
	i.End()
}

type streamReply struct {
	id   sourceId
	data inspect.Inspectable //data should be inspectable array for this message to be inspectable, same for streamRequest
}

type streamAck struct {
	id outputId
}

type upstreamStopped struct {
	id  sourceId
	err error
}

type downstreamStopped struct {
	id  outputId
	err error
}

type closeStream outputId
