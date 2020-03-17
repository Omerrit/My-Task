package queue

import (
	"gerrit-share.lan/go/common"
	"sync/atomic"
	"unsafe"
)

//because we all need another lock free queue
type Queue struct {
	head          unsafe.Pointer
	signalChannel common.SignalChannel
}

type QueueElement struct {
	Data interface{}
	Next *QueueElement
}

//idea stolen from CAF
func (q *Queue) closedTag() unsafe.Pointer {
	return unsafe.Pointer(q)
}

func (q *Queue) Push(data interface{}) bool {
	head := atomic.LoadPointer(&q.head)
	if head == q.closedTag() {
		return false
	}
	elem := &QueueElement{data, (*QueueElement)(head)}
	for !atomic.CompareAndSwapPointer(&q.head, head, unsafe.Pointer(elem)) {
		head = atomic.LoadPointer(&q.head)
		if head == q.closedTag() {
			return false
		}
		elem.Next = (*QueueElement)(head)
	}
	select {
	case q.signalChannel <- common.None{}:
	default:
	}
	return true
}

func (q *Queue) takeHead(newptr unsafe.Pointer) *QueueElement {
	head := atomic.LoadPointer(&q.head)
	if head == q.closedTag() || unsafe.Pointer(head) == newptr {
		return nil
	}
	for !atomic.CompareAndSwapPointer(&q.head, head, newptr) {
		head = atomic.LoadPointer(&q.head)
		if head == q.closedTag() || unsafe.Pointer(head) == newptr {
			return nil
		}
	}
	return (*QueueElement)(head)
}

func (q *Queue) TakeHead() *QueueElement {
	return q.takeHead(nil)
}

func (q *Queue) TakeHeadAndClose() *QueueElement {
	return q.takeHead(q.closedTag())
}

func (q *Queue) SignalChannel() common.OutSignalChannel {
	return q.signalChannel.ToOutput()
}

func NewQueue() *Queue {
	return &Queue{nil, make(common.SignalChannel, 1)}
}
