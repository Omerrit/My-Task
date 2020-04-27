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

type ids map[string]children

func (i *ids) acquireNewId(parent string) int {
	if *i == nil {
		*i = make(ids, 1)
	}
	children := (*i)[parent]
	max := children.acquireNewId()
	(*i)[parent] = children
	return max
}

func (i *ids) isRegistered(id string) bool {
	parent, numberAsString := splitId(id)
	children, ok := (*i)[parent]
	if !ok {
		return false
	}
	_, ok = children.ids[numberAsString]
	return ok
}

func (i *ids) deleteId(id string) error {
	parent, numberAsString := splitId(id)
	children, ok := (*i)[parent]
	if !ok {
		return fmt.Errorf("provided parent id is not registered")
	}
	children.deleteId(numberAsString)
	return nil
}

func (i *ids) reserveId(id string) error {
	parent, number := splitId(id)
	children := (*i)[parent]
	err := children.reserveId(number)
	if err != nil {
		return err
	}
	(*i)[parent] = children
	return nil
}

func (i *ids) restoreId(id string) error {
	if *i == nil {
		*i = make(ids, 1)
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

type TypedIds map[string]ids

func (t *TypedIds) AcquireNewId(objectType, parent string) string {
	if *t == nil {
		*t = make(TypedIds, 1)
	}
	objects := (*t)[objectType]
	newId := makeId(parent, objects.acquireNewId(parent))
	(*t)[objectType] = objects
	return newId
}

func (t *TypedIds) DeleteId(objectType, id string) error {
	objects, ok := (*t)[objectType]
	if !ok {
		return fmt.Errorf("provided type is not registered")
	}
	return objects.deleteId(id)
}

func (t *TypedIds) IsRegistered(objectType, id string) bool {
	objects, ok := (*t)[objectType]
	if !ok {
		return false
	}
	return objects.isRegistered(id)
}

func (t *TypedIds) ReserveId(objectType, id string) error {
	objects, ok := (*t)[objectType]
	if !ok {
		return fmt.Errorf("provided type is not registered")
	}
	return objects.reserveId(id)
}

func (t *TypedIds) RestoreId(objectType, id string) error {
	if *t == nil {
		*t = make(TypedIds, 1)
	}
	objects := (*t)[objectType]
	err := objects.restoreId(id)
	if err != nil {
		return err
	}
	(*t)[objectType] = objects
	return nil
}
