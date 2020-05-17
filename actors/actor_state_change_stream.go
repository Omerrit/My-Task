package actors

import (
	"gerrit-share.lan/go/inspect"
)

type ActorStateChangeStream struct {
	buffer            ActorsArray
	state             ActorSet
	startOffset       int64
	isAnyoneListening bool
}

func (a *ActorStateChangeStream) Add(actor ActorService) {
	a.state.Add(actor)
	if !a.isAnyoneListening {
		return
	}
	a.buffer = append(a.buffer, actor)
}

func (a *ActorStateChangeStream) Remove(actor ActorService) {
	a.state.Remove(actor)
}

func (a *ActorStateChangeStream) fillArray(array *ActorsArray, offset int64, maxLen int) (inspect.Inspectable, int64, error) {
	realLen := len(a.buffer) - int(offset-a.startOffset)
	if realLen > maxLen {
		realLen = maxLen
	}
	if realLen == 0 {
		return nil, offset, nil
	}
	array.SetLength(realLen)
	nextOffset := offset + int64(copy(*array, a.buffer[int(offset-a.startOffset):]))
	return array, nextOffset, nil
}

func (a *ActorStateChangeStream) FillData(data inspect.Inspectable, offset int64, maxLen int) (result inspect.Inspectable, nextOffset int64, err error) {
	if offset < a.startOffset {
		return data, offset, ErrOffsetOutOfRange
	}
	if maxLen == 0 {
		maxLen = DefaultMaxLen
	}
	if array, ok := data.(*ActorsArray); ok {
		return a.fillArray(array, offset, maxLen)
	} else if _, ok := data.(ActorService); ok {
		if int(offset-a.startOffset) > len(a.buffer) {
			return nil, offset, ErrOffsetOutOfRange
		}
		if int(offset-a.startOffset) == len(a.buffer) {
			return nil, offset, nil
		}
		return a.buffer[offset-a.startOffset], offset + 1, nil
	} else if data == nil {
		return a.fillArray(new(ActorsArray), offset, maxLen)
	}
	return data, offset, ErrWrongTypeRequested
}

func (a *ActorStateChangeStream) GetLatestState() (int64, DataSource) {
	array := make(ActorsArray, 0, len(a.state))
	for actor := range a.state {
		array = append(array, actor)
	}
	a.isAnyoneListening = true //start collecting history
	return a.startOffset + int64(len(a.buffer)), &StaticActorDataSource{Data: array}
}

//TODO: maybe make some parametrized criterion
func (a *ActorStateChangeStream) LastOffsetChanged(offset int64) {
	if len(a.buffer)/2 < int(offset-a.startOffset) {
		var buffer ActorsArray
		buffer.SetLength(len(a.buffer) - int(offset-a.startOffset))
		copy(buffer, a.buffer[int(offset-a.startOffset):])
		a.buffer = buffer
		a.startOffset = offset
	}
}

func (a *ActorStateChangeStream) NoMoreSubscribers() {
	a.buffer.SetLength(0)
	a.startOffset = 0
	a.isAnyoneListening = false
}

func init() {
	var _ StateChangeStream = (*ActorStateChangeStream)(nil)
}
