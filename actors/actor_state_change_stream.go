package actors

import (
	"gerrit-share.lan/go/inspect"
)

type ActorStateChangeStream struct {
	buffer            ActorsArray
	state             ActorSet
	startOffset       int
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

func (a *ActorStateChangeStream) fillArray(array *ActorsArray, offset int, maxLen int) (inspect.Inspectable, int, error) {
	realLen := len(a.buffer) - offset + a.startOffset
	if realLen > maxLen {
		realLen = maxLen
	}
	if realLen == 0 {
		return nil, offset, nil
	}
	array.SetLength(realLen)
	nextOffset := offset + copy(*array, a.buffer[(offset-a.startOffset):])
	return array, nextOffset, nil
}

func (a *ActorStateChangeStream) FillData(data inspect.Inspectable, offset int, maxLen int) (result inspect.Inspectable, nextOffset int, err error) {
	if offset < a.startOffset {
		return data, offset, ErrOffsetOutOfRange
	}
	if maxLen == 0 {
		maxLen = DefaultMaxLen
	}
	if array, ok := data.(*ActorsArray); ok {
		return a.fillArray(array, offset, maxLen)
	} else if _, ok := data.(ActorService); ok {
		if (offset - a.startOffset) > len(a.buffer) {
			return nil, offset, ErrOffsetOutOfRange
		}
		if (offset - a.startOffset) == len(a.buffer) {
			return nil, offset, nil
		}
		return a.buffer[offset-a.startOffset], offset + 1, nil
	} else {
		return a.fillArray(new(ActorsArray), offset, maxLen)
	}
	return data, offset, ErrWrongTypeRequested
}

func (a *ActorStateChangeStream) GetLatestState() (int, DataSource) {
	array := make(ActorsArray, 0, len(a.state))
	for actor := range a.state {
		array = append(array, actor)
	}
	a.isAnyoneListening = true //start collecting history
	return a.startOffset + len(a.buffer), &StaticActorDataSource{Data: array}
}

//TODO: maybe make some parametrized criterion
func (a *ActorStateChangeStream) LastOffsetChanged(offset int) {
	if len(a.buffer)/2 >= (offset - a.startOffset) {
		var buffer ActorsArray
		buffer.SetLength(len(a.buffer) + a.startOffset - offset)
		copy(buffer, a.buffer[offset-a.startOffset:])
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
