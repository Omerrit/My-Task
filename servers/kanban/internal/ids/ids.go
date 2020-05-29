package ids

import (
	"strconv"
)

type children struct {
	max int
	ids map[string]bool //id->is not empty
}

func (c *children) acquireNewId() int {
	c.max++
	_, ok := c.ids[strconv.Itoa(c.max)]
	if ok {
		return c.acquireNewId()
	}
	if c.ids == nil {
		c.ids = make(map[string]bool, 1)
	}
	c.ids[strconv.Itoa(c.max)] = false
	return c.max
}

func (c *children) deleteId(id string) bool {
	if c.ids[id] {
		return false
	}
	if !c.ids[id] {
		delete(c.ids, id)
	}
	number, err := strconv.Atoi(id)
	if err == nil && c.max == number {
		c.max--
	}
	return true
}

func (c *children) reserveId(id string) bool {
	_, ok := c.ids[id]
	if ok {
		return false
	}
	if c.ids == nil {
		c.ids = make(map[string]bool, 1)
	}
	c.ids[id] = false
	number, err := strconv.Atoi(id)
	if err == nil && c.max < number {
		c.max = number
	}
	return true
}

func (c *children) restoreId(id string) {
	c.reserveId(id)
	c.ids[id] = true
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

func (i *Ids) DeleteId(id string) bool {
	parent, numberAsString := splitId(id)
	children, ok := (*i)[parent]
	if !ok {
		return true
	}
	return children.deleteId(numberAsString)
}

func (i *Ids) ReserveId(id string) bool {
	if *i == nil {
		*i = make(Ids, 1)
	}
	parent, number := splitId(id)
	children := (*i)[parent]
	ok := children.reserveId(number)
	if !ok {
		return false
	}
	(*i)[parent] = children
	return true

}

func (i *Ids) RestoreId(id string) {
	if *i == nil {
		*i = make(Ids, 1)
	}
	parent, child := splitId(id)
	children := (*i)[parent]
	children.restoreId(child)
	(*i)[parent] = children
}
