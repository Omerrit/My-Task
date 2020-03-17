package actors

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

type ActorsArray []ActorService

const ActorsArrayName = packageName + ".array.actors"

func (a *ActorsArray) SetLength(length int) {
	if length > cap(*a) {
		*a = make(ActorsArray, length)
	} else {
		*a = (*a)[:length]
	}
}

func (a *ActorsArray) Inspect(inspector *inspect.GenericInspector) {
	i := inspector.Array(ActorsArrayName, ActorServiceName, "")
	{
		if i.IsReading() {
			a.SetLength(i.GetLength())
		}
		for index := range *a {
			InspectActorService(&((*a)[index]), i.Value())
		}
		i.End()
	}
}

func init() {
	inspectables.Register(ActorsArrayName, func() inspect.Inspectable { return new(ActorsArray) })
}
