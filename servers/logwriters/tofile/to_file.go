package tofile

import (
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/log"
	"gerrit-share.lan/go/actors/plugins/log/writerhelpers"
	"gerrit-share.lan/go/actors/starter"
	"gerrit-share.lan/go/utils/flags"
	"os"
)

const pluginName = "logfile_writer"

var fileFlag string

type logFileWriter struct {
	file *os.File
}

func newFileWriter() (*logFileWriter, error) {
	var err error
	fileWriter := &logFileWriter{}
	fileWriter.file, err = os.OpenFile(fileFlag, os.O_CREATE|os.O_APPEND|os.O_SYNC|os.O_RDWR, 0666)
	return fileWriter, err
}

func (fs *logFileWriter) Println(a ...interface{}) {
	if fs.file != nil {
		fmt.Fprintln(fs.file, a...)
	}
}

func (fs *logFileWriter) Printf(format string, a ...interface{}) {
	if fs.file != nil {
		fmt.Fprintf(fs.file, format, a...)
	}
}

func (fs *logFileWriter) Shutdown() error {
	if fs.file != nil {
		return fs.file.Close()
	}
	return nil
}

func init() {
	id := log.RegisterWriter()
	starter.SetCreator(pluginName, func(parent *actors.Actor, name string) (actors.ActorService, error) {
		if len(fileFlag) > 0 {
			fileWriter, err := newFileWriter()
			if err != nil {
				return nil, err
			}
			service := writerhelpers.NewPrintService(id, fileWriter, pluginName, 10)
			return parent.System().Spawn(service), nil
		}
		log.NotifyNotSubscribe(parent, id)
		return nil, actors.ErrNotGonnaHappen
	})

	starter.SetFlagInitializer(pluginName, func() {
		flags.StringFlag(&fileFlag, "logfile", "log file name")
	})
}
