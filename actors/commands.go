package actors

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

type GetInfo struct {
	DummyFiller
}

const GetInfoName = packageName + ".getinfo"

func (GetInfo) Inspect(inspector *inspect.GenericInspector) {
	inspector.Object(GetInfoName, "request info for all registered commands").End()
}

func init() {
	inspectables.Register(GetInfoName, func() inspect.Inspectable { return GetInfo{} })
}

type GetStatus struct {
	DummyFiller
}

const GetStatusName = packageName + ".getstatus"

func (GetStatus) Inspect(inspector *inspect.GenericInspector) {
	inspector.Object(GetStatusName, "request service dependent info and/or statistics").End()
}

func init() {
	inspectables.Register(GetStatusName, func() inspect.Inspectable { return GetStatus{} })
}

type Status struct {
	CommandTypes       int
	MessageTypes       int
	ProcessingCommands int
	AwaitingCommands   int
	StreamInputs       int
	StreamOutputs      int
}

const StatusName = packageName + ".status"

func (s *Status) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(StatusName, "status")
	{
		o.Int(&s.CommandTypes, "command_types", true, "number of commands this service can handle")
		o.Int(&s.MessageTypes, "message_types", true, "number of messages this service can handle")
		o.Int(&s.ProcessingCommands, "processing_commands", true, "number of commands this service is currently processing")
		o.Int(&s.AwaitingCommands, "awaiting_commands", true, "number of commands this service is currently awaiting")
		o.Int(&s.StreamInputs, "stream_inputs", true, "number of stream inputs this service is currently having")
		o.Int(&s.StreamOutputs, "stream_outputs", true, "number of stream outputs this service is currently having")
		o.End()
	}
}

func (s *Status) Visit(visitor ResponseVisitor) {
	visitor.Reply(s)
}

func init() {
	inspectables.Register(StatusName, func() inspect.Inspectable { return new(Status) })
}
