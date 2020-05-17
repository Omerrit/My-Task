package kanban

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

const packageName = "kanban"

type subscribe struct {
	actors.RequestStreamBase
	startId int64
}

func init() {
	inspectables.RegisterDescribed(packageName+".subscribe", func() inspect.Inspectable { return new(subscribe) }, "subscribe to broadcaster stream")
}

type message struct {
	key   string
	value string
}

const messageName = packageName + ".message"

func (m *message) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(messageName, "message")
	{
		objectInspector.String(&m.key, "key", true, "key of property")
		objectInspector.String(&m.value, "value", true, "value of property")
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(messageName, func() inspect.Inspectable { return new(message) })
}

type login struct {
	userName string
	password string
}

const loginName = packageName + ".login"

func (l *login) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(loginName, "login")
	{
		objectInspector.String(&l.userName, "user", true, "user name")
		objectInspector.String(&l.password, "password", true, "user password")
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(loginName, func() inspect.Inspectable { return new(login) })
}

type id struct {
	objectType string
	id         string
}

type newId id

const newIdName = packageName + ".newid"

func (n *newId) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(newIdName, "new id command")
	{
		objectInspector.String(&n.objectType, "type", true, "object type")
		objectInspector.String(&n.id, "parent", true, "parent id")
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(newIdName, func() inspect.Inspectable { return new(newId) })
}

type deleteId id

const deleteIdName = packageName + ".deleteid"

func (d *deleteId) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(deleteIdName, "id to delete")
	{
		objectInspector.String(&d.objectType, "type", true, "object type")
		objectInspector.String(&d.id, "id", true, "id to delete")
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(deleteIdName, func() inspect.Inspectable { return new(deleteId) })
}

type isIdRegistered id

const isIdRegisteredName = packageName + ".isidregistered"

func (r *isIdRegistered) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(isIdRegisteredName, "check if id is registered")
	{
		objectInspector.String(&r.objectType, "type", true, "object type")
		objectInspector.String(&r.id, "id", true, "id to check")
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(isIdRegisteredName, func() inspect.Inspectable { return new(isIdRegistered) })
}

type reserveId id

const reserveIdName = packageName + ".reserveid"

func (r *reserveId) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(reserveIdName, "reserve id")
	{
		objectInspector.String(&r.objectType, "type", true, "object type")
		objectInspector.String(&r.id, "id", true, "id to reserve")
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(reserveIdName, func() inspect.Inspectable { return new(reserveId) })
}
