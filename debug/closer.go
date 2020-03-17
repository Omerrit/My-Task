package debug

import (
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/interfaces"
	"io"
	"log"
)

type logCloser struct {
	closer io.Closer
}

func logError(err error, service io.Closer) error {
	if err != nil {
		log.Println(err)
	}
	return err
}

// TODO: better info
func (l logCloser) Close() error {
	return logError(l.closer.Close(), l.closer)
}

type logClosableService struct {
	closableService interfaces.ClosableService
}

func (l logClosableService) Close() error {
	return logError(l.closableService.Close(), l.closableService)
}

func (l logClosableService) DoneChannel() common.OutSignalChannel {
	return l.closableService.DoneChannel()
}

func (l logClosableService) Shutdown() {
	l.closableService.Shutdown()
}

func LogErrorOnClose(closer io.Closer) io.Closer {
	if closableService, ok := closer.(interfaces.ClosableService); ok {
		return logClosableService{closableService}
	}
	return logCloser{closer}
}
