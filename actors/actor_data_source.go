package actors

import (
	"gerrit-share.lan/go/inspect"
)

//deafult length for stream source when maxLen==0
const DefaultMaxLen = 10

type StaticActorDataSource struct {
	Data          ActorsArray
	startOffset   int
	currentOffset int
}

func (s *StaticActorDataSource) Acknowledged() {
	s.startOffset = s.currentOffset
	if s.currentOffset == len(s.Data) {
		s.Data = nil
	}
}

func (s *StaticActorDataSource) fillArray(array *ActorsArray, maxLen int) inspect.Inspectable {
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

func (s *StaticActorDataSource) FillData(data inspect.Inspectable, maxLen int) (inspect.Inspectable, error) {
	if len(s.Data) == 0 {
		return nil, nil
	}
	if maxLen == 0 {
		maxLen = DefaultMaxLen
	}
	if array, ok := data.(*ActorsArray); ok {
		return s.fillArray(array, maxLen), nil
	} else if value, ok := data.(ActorService); ok {
		//ignore maxLen
		if s.startOffset >= len(s.Data) {
			return nil, nil
		}
		value = s.Data[s.startOffset]
		s.currentOffset = s.startOffset + 1
		return value, nil
	} else if data == nil {
		return s.fillArray(new(ActorsArray), maxLen), nil
	}
	return nil, ErrWrongTypeRequested
}

func init() {
	var _ DataSource = (*StaticActorDataSource)(nil)
}
