package log

import (
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectwrappers"
)

func toInspectable(data interface{}) inspect.Inspectable {
	switch v := data.(type) {
	case nil:
		return nil
	case inspect.Inspectable:
		return v
	case fmt.Stringer:
		return inspectwrappers.NewStringValue(v.String())
	default:
		return inspectwrappers.NewStringValue(fmt.Sprintf("%#v", v))
	}
}

func writeMessage(actor actors.ActorCompatible, severity Severity, stack errors.CallStack, source string, message string, extra interface{}, sourceInfo interface{}) {
	msg := &Message{
		Severity:   severity,
		Stack:      stack,
		Source:     source,
		Message:    message,
		Extra:      toInspectable(extra),
		SourceInfo: toInspectable(sourceInfo)}
	base := actor.GetBase()
	base.SendMessage(base.System().GetPluginActor(pluginName), msg)
}

func writeSimpleMessage(actor actors.ActorCompatible, severity Severity, stack errors.CallStack, message string) {
	writeMessage(actor, severity, stack, actor.GetBase().GetCommandInfo().Name, message, nil, nil)
}

func writeError(actor actors.ActorCompatible, stackTrace errors.CallStack, err error, severity Severity, source string, sourceInfo interface{}) {
	if err == nil {
		return
	}
	switch val := err.(type) {
	case errors.ErrorArray:
		writeMessage(actor, severity, stackTrace, source, val.Error(), val.ToError(), sourceInfo)
	case *errors.GoPanic:
		writeError(actor, stackTrace, val.Unwrap(), SeverityCrash, source, sourceInfo)
	case errors.UnknownStackTraceError:
		writeMessage(actor, severity, stackTrace, source, val.Error(), val.Source(), sourceInfo)
	case errors.StackTraceError:
		if val.Unwrap() == nil {
			writeMessage(actor, severity, stackTrace, source, val.Error(), nil, sourceInfo)
			return
		}
		writeError(actor, stackTrace, val.Unwrap(), severity, source, sourceInfo)
	default:
		writeSimpleMessage(actor, severity, stackTrace, val.Error())
	}
}

func writeErrorWrapper(actor actors.ActorCompatible, err error, severity Severity, stack errors.CallStack,
	source string, sourceInfo interface{}) {
	if err == nil {
		return
	}
	if len(source) == 0 {
		source = actor.GetBase().GetCommandInfo().Name
	}

	switch val := err.(type) {
	case errors.ErrorArray:
		for _, errorItem := range val {
			errStackTrace, ok := errorItem.(errors.StackTraceError)
			if ok {
				writeError(actor, errStackTrace.StackTrace(), errStackTrace, severity, source, sourceInfo)
				return
			}
			writeError(actor, stack, errorItem, severity, source, sourceInfo)
		}
	case errors.StackTraceError:
		writeError(actor, val.StackTrace(), val, severity, source, sourceInfo)
	default:
		writeError(actor, stack, err, severity, source, sourceInfo)
	}
}

func writeErrorNoInfo(actor actors.ActorCompatible, err error, severity Severity, stack errors.CallStack) {
	writeErrorWrapper(actor, err, severity, stack, "", nil)
}
