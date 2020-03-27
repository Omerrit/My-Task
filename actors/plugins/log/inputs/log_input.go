package inputs

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/log"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"strings"
	"time"
)

type Printer interface {
	Println(a ...interface{})
	Printf(format string, a ...interface{})
}

type logInput struct {
	actors.StreamInputBase
	printer Printer
	maxLen  int
}

func NewLogInput(maxLen int, printer Printer) *logInput {
	return &logInput{
		printer: printer,
		maxLen:  maxLen,
	}
}

func (dw *logInput) Process(data inspect.Inspectable) error {
	msgs := data.(*log.Messages)
	for _, msg := range *msgs {
		err := dw.processMessage(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dw *logInput) RequestNext() {
	dw.StreamInputBase.RequestData(new(log.Messages), dw.maxLen)
}

func (dw *logInput) processMessage(message log.Message) error {
	if dw.printer == nil {
		return nil
	}
	if message.Severity > log.GetSeverity() {
		return nil
	}
	now := time.Now()
	dw.printer.Printf("%v-%v-%v %v:%v:%v %v: %v\n", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second(), strings.ToUpper(log.GetSeverityString(message.Severity)), strings.TrimSuffix(message.Message, "\n"))
	if len(message.Source) > 0 {
		dw.printer.Println("source:", message.Source)
	}
	if message.SourceInfo != nil {
		dw.printer.Printf("source info: %#v\n", message.SourceInfo)
	}
	if message.Severity == log.SeverityDebug || message.Severity <= log.SeverityError {
		message.Stack.CutTop(func(frame errors.StackFrame) bool {
			if frame.Package == "runtime" {
				return true
			}
			return false
		})
		dw.printer.Println(message.Stack)
	}
	if log.GetVerbosity() >= log.VerbosityHigh && message.Extra != nil {
		dw.printer.Println("extra:", message.Extra)
	}
	return nil
}
