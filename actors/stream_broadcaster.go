package actors

import (
	"container/heap"
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/inspect"
)

type StateChangeBroadcasterOutput struct {
	StreamOutputBase
	broadcaster           StateBroadcaster
	myOffset              int64
	nextOffset            int64
	stateSource           DataSource
	isInitialDataFinished bool
}

func (out *StateChangeBroadcasterOutput) DataSource() DataSource {
	return out.stateSource
}

func (out *StateChangeBroadcasterOutput) Acknowledged() {
	if !out.isInitialDataFinished {
		out.stateSource.Acknowledged()
		return
	}
	out.broadcaster.offsetChanged(out, out.nextOffset)
	out.myOffset = out.nextOffset
}

func (out *StateChangeBroadcasterOutput) FillData(data inspect.Inspectable, maxLen int) (inspect.Inspectable, error) {
	var result inspect.Inspectable
	var err error
	if !out.isInitialDataFinished {
		if out.stateSource != nil {
			result, err = out.stateSource.FillData(data, maxLen)
			if err != nil {
				return nil, err
			}
			if result != nil {
				return result, nil
			}
			if dynamic, ok := out.stateSource.(DynamicDataSource); ok {
				if dynamic.RequestNext(out.FlushLater, out.myOffset) {
					return out.stateSource.FillData(data, maxLen)
				}
			}
		}
		out.isInitialDataFinished = true
	}
	result, out.nextOffset, err = out.broadcaster.fillData(out, data, out.myOffset, maxLen)
	return result, err
}

func (out *StateChangeBroadcasterOutput) resetOffset(newZero int64) {
	out.myOffset -= newZero
	out.nextOffset -= newZero
}

func (out *StateChangeBroadcasterOutput) restart() {
	out.isInitialDataFinished = false
	out.FlushLater()
}

type offset struct {
	index int
	out   *StateChangeBroadcasterOutput
}

type offsets []*offset

func (o offsets) Len() int {
	return len(o)
}

func (o offsets) Less(i int, j int) bool {
	return o[i].out.myOffset < o[j].out.myOffset
}

func (o offsets) Swap(i int, j int) {
	o[i], o[j] = o[j], o[i]
	o[i].index = i
	o[j].index = j
}

func (o *offsets) Push(value interface{}) {
	value.(*offset).index = len(*o)
	*o = append(*o, value.(*offset))
}

func (o *offsets) Pop() interface{} {
	length := len(*o)
	result := (*o)[length-1]
	*o = (*o)[0:(length - 1)]
	return result
}

type stateChangeBroadcasterImpl struct {
	outputs                    map[*StateChangeBroadcasterOutput]*offset
	readyOutputs               map[*StateChangeBroadcasterOutput]common.None
	offsets                    offsets
	stream                     StateChangeStream
	dropPolicy                 DropPolicy
	shouldCloseWhenActorCloses bool
}

func (b *stateChangeBroadcasterImpl) fillData(out *StateChangeBroadcasterOutput, data inspect.Inspectable, offset int64, maxLen int) (inspect.Inspectable, int64, error) {
	outData, outOffset, err := b.stream.FillData(data, offset, maxLen)
	if err == nil && outData == nil {
		b.readyOutputs[out] = common.None{}
	}
	return outData, outOffset, err
}

func (b *stateChangeBroadcasterImpl) offsetChanged(out *StateChangeBroadcasterOutput, newOffset int64) {
	off := b.outputs[out]
	if off.index == 0 {
		heap.Fix(&b.offsets, off.index)
		b.stream.LastOffsetChanged(b.offsets[0].out.myOffset)
		return
	}
	heap.Fix(&b.offsets, off.index)
}

func (b *stateChangeBroadcasterImpl) removeOut(out *StateChangeBroadcasterOutput) {
	off := b.outputs[out]
	heap.Remove(&b.offsets, off.index)
	delete(b.outputs, out)
	delete(b.readyOutputs, out)
	if off.index == 0 {
		if b.offsets.Len() > 0 {
			b.stream.LastOffsetChanged(b.offsets[0].out.myOffset)
		} else {
			b.stream.NoMoreSubscribers()
		}
	}
}

func (b *stateChangeBroadcasterImpl) restart(out *StateChangeBroadcasterOutput) {
	off := b.outputs[out]
	out.myOffset, out.stateSource = b.stream.GetLatestState()
	if off.index == 0 {
		heap.Fix(&b.offsets, off.index)
		b.stream.LastOffsetChanged(b.offsets[0].out.myOffset)
		out.restart()
		return
	}
	heap.Fix(&b.offsets, off.index)
	out.restart()
}

func (b *stateChangeBroadcasterImpl) AddOutput() *StateChangeBroadcasterOutput {
	out := new(StateChangeBroadcasterOutput)
	out.broadcaster = b
	out.myOffset, out.stateSource = b.stream.GetLatestState()
	off := &offset{0, out}
	b.outputs[out] = off
	heap.Push(&b.offsets, off)
	out.OnClose(func(error) {
		b.removeOut(out)
	})
	if b.shouldCloseWhenActorCloses {
		out.CloseWhenActorCloses()
	}
	out.restart()
	return out
}

func (b *stateChangeBroadcasterImpl) applyDropPolicy() {
	for b.offsets.Len() > 0 {
		if !b.dropPolicy(b.offsets[0].out.myOffset) {
			return
		}
		b.offsets[0].out.CloseStream(ErrStreamConsumerSlow)
	}
}

func (b *stateChangeBroadcasterImpl) SetDropPolicy(policy DropPolicy) {
	b.dropPolicy = policy
}

func (b *stateChangeBroadcasterImpl) NewDataAvailable() {
	if b.dropPolicy != nil {
		b.applyDropPolicy()
	}
	for out := range b.readyOutputs {
		out.FlushLater()
	}
	for out := range b.readyOutputs {
		delete(b.readyOutputs, out)
	}
}

func (b *stateChangeBroadcasterImpl) Close(err error) {
	for out := range b.outputs {
		out.CloseStream(err)
	}
}

func (b *stateChangeBroadcasterImpl) CloseWhenActorCloses() {
	b.shouldCloseWhenActorCloses = true
}

func NewBroadcaster(stream StateChangeStream) StateBroadcaster {
	return &stateChangeBroadcasterImpl{make(map[*StateChangeBroadcasterOutput]*offset), make(map[*StateChangeBroadcasterOutput]common.None), nil, stream, nil, false}
}

type DropPolicy func(offset int64) bool

type StateBroadcaster interface {
	fillData(out *StateChangeBroadcasterOutput, data inspect.Inspectable, offset int64, maxLen int) (inspect.Inspectable, int64, error)
	offsetChanged(out *StateChangeBroadcasterOutput, newOffset int64)
	AddOutput() *StateChangeBroadcasterOutput
	NewDataAvailable()
	SetDropPolicy(DropPolicy)
	Close(err error)
	CloseWhenActorCloses()
}
