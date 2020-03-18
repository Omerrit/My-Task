package writerhelpers

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/log"
	"gerrit-share.lan/go/actors/plugins/log/inputs"
	"gerrit-share.lan/go/interfaces"
)

type ShutdownablePrinter interface {
	inputs.Printer
	interfaces.Shutdownable
}

type PrintService struct {
	actors.Actor
	id      log.WriterId
	printer ShutdownablePrinter
	name    string
	maxLen  int
}

func NewPrintService(id log.WriterId, printer ShutdownablePrinter, name string, maxLen int) *PrintService {
	return &PrintService{
		id:      id,
		printer: printer,
		name:    name,
		maxLen:  maxLen,
	}
}

func (d *PrintService) MakeBehaviour() actors.Behaviour {
	d.printer.Println(d.name, "started")
	log.SubscribeForMessages(d, d.id, inputs.NewLogInput(d.maxLen, d.printer))
	return actors.Behaviour{}
}

func (d *PrintService) Shutdown() error {
	d.printer.Println(d.name, "shut down")
	return d.printer.Shutdown()
}
