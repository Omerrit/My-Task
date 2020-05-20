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

type saveFile struct {
	lastModified int
	user         string
	id           string
}

const saveFileName = packageName + ".savefile"

func (s *saveFile) Embed(o *inspect.ObjectInspector) {
	o.Int(&s.lastModified, "modified", true, "time of last modification of file")
	o.String(&s.user, "user", true, "user")
	o.String(&s.id, "id", true, "id")
}

func (s *saveFile) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(saveFileName, "save file command")
	{
		s.Embed(objectInspector)
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(saveFileName, func() inspect.Inspectable { return new(saveFile) })
}

type saveMsgToKafka struct {
	value string
	key   string
}

const saveMsgToKafkaName = packageName + ".savemsgtokafka"

func (s *saveMsgToKafka) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(saveMsgToKafkaName, "save msg to kafka")
	{
		objectInspector.String(&s.value, "value", false, "msg value")
		objectInspector.String(&s.key, "key", true, "msg key")
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(saveMsgToKafkaName, func() inspect.Inspectable { return new(saveMsgToKafka) })
}

type saveMsgsToKafka []saveMsgToKafka

const saveMsgsToKafkaName = packageName + ".savemsgstokafka"

func (s *saveMsgsToKafka) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(saveMsgsToKafkaName, saveMsgToKafkaName, "readable/writable bigInt array")
	{
		if !arrayInspector.IsReading() {
			arrayInspector.SetLength(len(*s))
		} else {
			length := arrayInspector.GetLength()
			if cap(*s) > length {
				*s = (*s)[:length]
			} else {
				*s = make([]saveMsgToKafka, length)
			}
		}
		for index := range *s {
			(*s)[index].Inspect(arrayInspector.Value())
		}
		arrayInspector.End()
	}
}

func init() {
	inspectables.Register(saveMsgsToKafkaName, func() inspect.Inspectable { return new(saveMsgsToKafka) })
}
