package kafka

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"log"
)

type consumerOutput struct {
	actors.StreamOutputBase
	Messages      Messages
	currentOffset int
}

func (c *consumerOutput) Acknowledged() {
	c.Messages = c.Messages[:copy(c.Messages, c.Messages[c.currentOffset:])]
	c.currentOffset = 0
}

func (c *consumerOutput) fillArray(array *Messages, maxLen int) inspect.Inspectable {
	realLen := len(c.Messages)
	if realLen > maxLen {
		realLen = maxLen
	}
	if realLen == 0 {
		return nil
	}
	array.SetLength(realLen)
	copy(*array, c.Messages)
	if realLen > c.currentOffset {
		c.currentOffset = realLen
	}
	return array
}

func (c *consumerOutput) FillData(data inspect.Inspectable, maxLen int) (inspect.Inspectable, error) {
	if c.currentOffset != 0 && c.GetActor() != nil {
		log.Printf("%p [%s]: double fill\n", c.GetActor().Service(), c.GetActor().Service().Name())
	} else if c.currentOffset != 0 {
		log.Println("double fill")
	}
	if len(c.Messages) == 0 {
		return nil, nil
	}
	if maxLen == 0 {
		maxLen = actors.DefaultMaxLen
	}
	if array, ok := data.(*Messages); ok {
		return c.fillArray(array, maxLen), nil
	} else if value, ok := data.(*Message); ok {
		if len(c.Messages) == 0 {
			return nil, nil
		}
		value = c.Messages[0]
		if c.currentOffset == 0 {
			c.currentOffset = 1
		}
		return value, nil
	} else if data == nil {
		return c.fillArray(new(Messages), maxLen), nil
	}
	return nil, actors.ErrWrongTypeRequested
}

func init() {
	var _ actors.DataSource = (*consumerOutput)(nil)
}
