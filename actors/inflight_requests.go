package actors

import (
	"gerrit-share.lan/go/interfaces"
)

type requestInfo struct {
	destination ActorService
	processor   interfaces.ReplyProcessor
}

type inflightRequests map[commandId]requestInfo

func (r *inflightRequests) Add(id commandId, processor interfaces.ReplyProcessor) {
	if *r == nil {
		*r = make(inflightRequests, 1)
	}
	(*r)[id] = requestInfo{nil, processor}
}

func (r *inflightRequests) SetDestination(id commandId, destination ActorService) {
	info, ok := (*r)[id]
	if ok {
		info.destination = destination
		(*r)[id] = info
	}
}

func (r *inflightRequests) Contains(id commandId) bool {
	_, ok := (*r)[id]
	return ok
}

func (r *inflightRequests) callProcessor(id commandId, caller func(interfaces.ReplyProcessor)) {
	info, ok := (*r)[id]
	if info.processor == nil {
		if ok {
			delete(*r, id)
		}
		return
	}
	delete(*r, id)
	caller(info.processor)
}

func (r *inflightRequests) Process(id commandId, reply interface{}) {
	r.callProcessor(id, func(processor interfaces.ReplyProcessor) {
		processor.Process(reply)
	})
}

func (r *inflightRequests) Error(id commandId, err error) {
	r.callProcessor(id, func(processor interfaces.ReplyProcessor) {
		processor.Error(err)
	})
}

func (r inflightRequests) IsEmpty() bool {
	return len(r) == 0
}

func (r *inflightRequests) Clear() {
	*r = make(inflightRequests, 1)
}
