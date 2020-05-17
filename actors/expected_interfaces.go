package actors

import (
	"gerrit-share.lan/go/inspect"
)

type (
	Response interface {
		Visit(ResponseVisitor)
	}

	StateChangeStream interface {
		FillData(data inspect.Inspectable, offset int64, maxLen int) (result inspect.Inspectable, nextOffset int64, err error)
		LastOffsetChanged(offset int64)
		GetLatestState() (lastOffset int64, lastStateSource DataSource)
		NoMoreSubscribers()
	}

	DataSource interface {
		Acknowledged()
		FillData(data inspect.Inspectable, maxLen int) (inspect.Inspectable, error)
	}
	DynamicDataSource interface {
		DataSource
		RequestNext(onDataAvailable func(), maxOffset int64) bool
	}
)

type DummyFiller struct {
}

func (DummyFiller) Fill() {

}
