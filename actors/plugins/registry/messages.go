package registry

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

const packageName = "actplugins.registry"

type actorName struct {
	name string
}

func (a *actorName) Fill() {
	a.name = ""
}

func (a *actorName) Embed(inspector *inspect.ObjectInspector) {
	inspector.String(&a.name, "name", true, "actor name")
}

func (a *actorName) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object("", "")
	a.Embed(o)
	o.End()
}

type registerMe struct {
	actorName
}

type registerOther struct {
	actorName
	actor actors.ActorService
}

const registerOtherName = packageName + ".register2"

func (r *registerOther) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(registerOtherName, "")
	{
		r.Embed(o)
		actors.InspectActorService(&r.actor, o.Value("actor", true, "actor to register"))
		o.End()
	}
}

type getActor struct {
	actorName
}

type waitActor struct {
	actorName
}

type willNotRegister struct {
	actorName
}

func init() {
	inspectables.RegisterDescribed(packageName+".register", func() inspect.Inspectable { return new(registerMe) },
		"register sender")
	inspectables.RegisterDescribed(registerOtherName, func() inspect.Inspectable { return new(registerOther) },
		"register any actor")
	inspectables.RegisterDescribed(packageName+".getactor", func() inspect.Inspectable { return new(getActor) },
		"get actor, return immediately if it's not available")
	inspectables.RegisterDescribed(packageName+".waitactor", func() inspect.Inspectable { return new(waitActor) },
		"get actor, wait for it to be registered if there is currently none")
	inspectables.RegisterDescribed(packageName+".will_not_register", func() inspect.Inspectable { return new(willNotRegister) },
		`mark name as unregisterable, registration with this 
		name will fail and all waiting for this name will get error immediately.
		Fails if some actor registered under that name before this call`)
}

type subscribe struct {
	actors.RequestStreamBase
}

func init() {
	inspectables.RegisterDescribed(packageName+".subscribe", func() inspect.Inspectable { return new(subscribe) }, "start registration info stream")
}

type Info struct {
	Name  string
	Actor actors.ActorService
}

const InfoName = packageName + ".info"

func (i *Info) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(InfoName, "registered actor")
	{
		o.String(&i.Name, "name", true, "registered actor name")
		actors.InspectActorService(&i.Actor, o.Value("actor", true, ""))
		o.End()
	}
}

type InfoArray []Info

const InfoArrayName = packageName + ".array.info"

func (a *InfoArray) SetLength(length int) {
	if cap(*a) < length {
		*a = make(InfoArray, length)
	} else {
		*a = (*a)[:length]
	}

}

func (a *InfoArray) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Array(InfoArrayName, InfoName, "registered actors")
	{
		if o.IsReading() {
			a.SetLength(o.GetLength())
		} else {
			o.SetLength(len(*a))
		}
		for index := range *a {
			(*a)[index].Inspect(o.Value())
		}
		o.End()
	}
}

func init() {
	inspectables.Register(InfoName, func() inspect.Inspectable { return new(Info) })
	inspectables.Register(InfoArrayName, func() inspect.Inspectable { return new(InfoArray) })
}
