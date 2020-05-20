package ids

import (
	"fmt"
	"gerrit-share.lan/go/common"
	"strconv"
)

type children struct {
	max int
	ids map[string]common.None
}

func (c *children) acquireNewId() int {
	c.max++
	_, ok := c.ids[strconv.Itoa(c.max)]
	if ok {
		return c.acquireNewId()
	}
	if c.ids == nil {
		c.ids = make(map[string]common.None, 1)
	}
	c.ids[strconv.Itoa(c.max)] = common.None{}
	return c.max
}

func (c *children) deleteId(id string) {
	delete(c.ids, id)
	number, err := strconv.Atoi(id)
	if err == nil && c.max == number {
		c.max--
	}
}

func (c *children) reserveId(id string) error {
	_, ok := c.ids[id]
	if ok {
		return fmt.Errorf("id is already in use for provided type")
	}
	if c.ids == nil {
		c.ids = make(map[string]common.None, 1)
	}
	c.ids[id] = common.None{}
	number, err := strconv.Atoi(id)
	if err == nil && c.max < number {
		c.max = number
	}
	return nil
}

type Ids map[string]children

func (i *Ids) AcquireNewId(parent string) string {
	if *i == nil {
		*i = make(Ids, 1)
	}
	children := (*i)[parent]
	max := children.acquireNewId()
	(*i)[parent] = children
	return makeId(parent, max)
}

func (i *Ids) IsRegistered(id string) bool {
	parent, numberAsString := splitId(id)
	children, ok := (*i)[parent]
	if !ok {
		return false
	}
	_, ok = children.ids[numberAsString]
	return ok
}

func (i *Ids) DeleteId(id string) error {
	parent, numberAsString := splitId(id)
	children, ok := (*i)[parent]
	if !ok {
		return fmt.Errorf("provided parent id is not registered")
	}
	children.deleteId(numberAsString)
	return nil
}

func (i *Ids) RestoreId(id string) error {
	if *i == nil {
		*i = make(Ids, 1)
	}
	parent, number := splitId(id)
	children := (*i)[parent]
	err := children.reserveId(number)
	if err != nil {
		return err
	}
	(*i)[parent] = children
	return nil
}
