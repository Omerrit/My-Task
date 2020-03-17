package registry

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
)

type staticInfoDataSource struct {
	Data          InfoArray
	startOffset   int
	currentOffset int
}

func (s *staticInfoDataSource) Acknowledged() {
	s.startOffset = s.currentOffset
	if s.currentOffset == len(s.Data) {
		s.Data = nil
	}
}

func (s *staticInfoDataSource) fillArray(array *InfoArray, maxLen int) inspect.Inspectable {
	realLen := len(s.Data) - s.startOffset
	if realLen > maxLen {
		realLen = maxLen
	}
	if realLen == 0 {
		return nil
	}
	array.SetLength(realLen)
	s.currentOffset = s.startOffset + copy(*array, s.Data[s.startOffset:])
	return array
}

func (s *staticInfoDataSource) FillData(data inspect.Inspectable, maxLen int) (inspect.Inspectable, error) {
	if len(s.Data) == 0 {
		return nil, nil
	}
	if maxLen == 0 {
		maxLen = actors.DefaultMaxLen
	}
	if array, ok := data.(*InfoArray); ok {
		return s.fillArray(array, maxLen), nil
	} else if value, ok := data.(*Info); ok {
		//ignore maxLen
		if s.startOffset >= len(s.Data) {
			return nil, nil
		}
		*value = s.Data[s.startOffset]
		s.currentOffset = s.startOffset + 1
		return value, nil
	} else if data == nil {
		return s.fillArray(new(InfoArray), maxLen), nil
	}
	return nil, actors.ErrWrongTypeRequested
}

func init() {
	var _ actors.DataSource = (*staticInfoDataSource)(nil)
}
