package actors

import (
	"github.com/ef-ds/deque"
)

type commandPopper interface {
	pop() resumedCommand
	len() int
}

type Command struct {
	cmd resumedCommand
}

func (c *Command) pop() resumedCommand {
	cmd := c.cmd
	c.cmd.invalidate()
	return cmd
}

func (c *Command) len() int {
	if !c.cmd.isValid() {
		return 0
	}
	return 1
}

func (c *Command) invalidate() {
	c.cmd.invalidate()
}

type CommandQueue struct {
	queue deque.Deque
	ready promiseIdSet
}

func (c *CommandQueue) push(cmd resumedCommand) {
	c.queue.PushFront(cmd)
	c.ready.Add(cmd.promiseId)
}

func (c *CommandQueue) pop() resumedCommand {
	var cmd resumedCommand
	var item interface{}
	var ok bool
	for c.queue.Len() > 0 {
		item, ok = c.queue.PopBack()
		if !ok {
			return resumedCommand{}
		}
		cmd = item.(resumedCommand)
		if !c.ready.Contains(cmd.promiseId) {
			continue
		}
		c.ready.Delete(cmd.promiseId)
		return cmd
	}
	return resumedCommand{}
}

func (c *CommandQueue) len() int {
	return c.queue.Len()
}

func (c *CommandQueue) rebuild() {
	var queue deque.Deque
	var item interface{}
	var ok bool
	var id promiseId
	for c.queue.Len() > 0 {
		item, ok = c.queue.PopBack()
		if !ok {
			break
		}
		id = item.(commandMessage).promiseId
		if !c.ready.Contains(id) {
			continue
		}
		queue.PushFront(item)
	}
	c.queue = queue
}

func (c *CommandQueue) addCanceled(id promiseId) {
	c.ready.Delete(id)
	//just some random criterion
	if len(c.ready) < (c.queue.Len() / 2) {
		c.rebuild()
	}
}

type resumedCommand struct {
	commandMessage
	filterIndex int
}

type reissuedQueue struct {
	queue deque.Deque
}

func (r *reissuedQueue) push(cmd resumedCommand) {
	r.queue.PushFront(cmd)
}

func (r *reissuedQueue) pushOne(queue commandPopper) {
	if queue.len() == 0 {
		return
	}
	r.queue.PushFront(queue.pop())
}

func (r *reissuedQueue) pushAll(queue commandPopper) {
	for queue.len() > 0 {
		r.queue.PushFront(queue.pop())
	}
}

func (r *reissuedQueue) len() int {
	return r.queue.Len()
}

func (r *reissuedQueue) pop() resumedCommand {
	cmd, ok := r.queue.PopBack()
	if !ok {
		return resumedCommand{}
	}
	return cmd.(resumedCommand)
}

//current command may be saved with push and then at later time one or more commands could be resend to service queue
//to be reprocessed in natural order
//note: need to activate cancellation mechanism when saving command and check if command was already canceled when processing
//that 'command to be reprocessed'

//note: there may be several queues, user should store them Actor should only provide functionality to store/resend
//reply visitor also should note that it's inside reprocessed command, it should replace cancel callback only if key is present in active promises map

//also there should be a way to cancel entire queue
