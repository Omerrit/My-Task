package consolelogwriter

import (
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/log"
	"gerrit-share.lan/go/actors/plugins/log/writerhelpers"
	"gerrit-share.lan/go/common"
)

const pluginName = "console_writer"

type consoleWriter common.None

func (c *consoleWriter) Println(a ...interface{}) {
	fmt.Println(a...)
}

func (c *consoleWriter) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (c *consoleWriter) Shutdown() error {
	return nil
}

func init() {
	id := log.RegisterWriter()
	actors.AddPlugin(pluginName, func() actors.BehavioralActor {
		return writerhelpers.NewPrintService(id, &consoleWriter{}, pluginName, 10)
	})
}
