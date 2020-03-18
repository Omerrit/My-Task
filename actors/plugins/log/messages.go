package log

import (
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

const packageName = "actplugins.log"

type Message struct {
	Severity   Severity
	Stack      errors.CallStack
	Source     string
	Message    string
	Extra      inspect.Inspectable
	SourceInfo inspect.Inspectable
}

const MessageName = packageName + ".message"

func (m *Message) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(MessageName, "logger message")
	{
		m.Severity.Embed(objectInspector)
		objectInspector.String(&m.Source, "source", false, "source name")
		objectInspector.String(&m.Message, "message", true, "message to log")
		genericInspector := objectInspector.Value("stacktrace", true, "stack trace")
		arrayInspector := genericInspector.Array(packageName+".stacktrace", "", "stack trace")
		{
			if !arrayInspector.IsReading() {
				arrayInspector.SetLength(len(m.Stack))
			} else {
				m.Stack = make(errors.CallStack, arrayInspector.GetLength())
			}
			for index := range m.Stack {
				genericInspector := arrayInspector.Value()

				objectInspector := genericInspector.Object("", "")
				{
					objectInspector.String(&m.Stack[index].File, "file", true, "file name")
					objectInspector.String(&m.Stack[index].Function, "function", true, "function name")
					objectInspector.String(&m.Stack[index].Package, "package", true, "package name")
					objectInspector.Int(&m.Stack[index].Line, "line", true, "line number")
					objectInspector.Int(&m.Stack[index].Offset, "file", true, "file name")
					objectInspector.End()
				}
			}
			arrayInspector.End()
		}
		if !i.IsReading() {
			if m.Extra != nil {
				m.Extra.Inspect(objectInspector.Value("extra", false, "extra information"))
			}
			if m.SourceInfo != nil {
				m.SourceInfo.Inspect(objectInspector.Value("sourceinfo", false, "source information"))
			}
		} else {
			m.Extra = nil
			m.SourceInfo = nil
		}
		objectInspector.End()
	}
}

func init() {
	inspectables.RegisterDescribed(MessageName, func() inspect.Inspectable { return new(Message) }, "new logger message")
}

type Messages []Message

const MessagesName = packageName + ".messages"

func (m *Messages) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(MessagesName, MessageName, "array of logger messages")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*m))
		} else {
			m.SetLength(arrayInspector.GetLength())
		}
		for index := range *m {
			(*m)[index].Inspect(arrayInspector.Value())
		}
		arrayInspector.End()
	}
}

func (m *Messages) SetLength(length int) {
	if cap(*m) < length {
		*m = make(Messages, length)
	} else {
		*m = (*m)[:length]
	}
}

func init() {
	inspectables.Register(MessagesName, func() inspect.Inspectable { return new(Messages) })
}
