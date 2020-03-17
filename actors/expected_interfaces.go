package actors

import (
	"gerrit-share.lan/go/inspect"
)

type (
	Response interface {
		Visit(ResponseVisitor)
	}

	StateChangeStream interface {
		FillData(data inspect.Inspectable, offset int, maxLen int) (result inspect.Inspectable, nextOffset int, err error)
		LastOffsetChanged(offset int)
		GetLatestState() (lastOffset int, lastStateSource DataSource)
		NoMoreSubscribers()
	}

	DataSource interface {
		Acknowledged()
		FillData(data inspect.Inspectable, maxLen int) (inspect.Inspectable, error)
	}
)

type DummyFiller struct {
}

func (DummyFiller) Fill() {

}
