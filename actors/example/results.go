package example

//result types go here
import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

//getStuff command returns this
type stuff struct {
	name  string
	value int
}

const stuffName = packageName + ".stuff"

func (s *stuff) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(stuffName, "example result,random stuff")
	{
		o.String(&s.name, "name", true, "some name")
		o.Int(&s.value, "value", true, "some value")
		o.End()
	}
}

//make result type a simple reply so it can be simply returned by command processor
func (s *stuff) Visit(visitor actors.ResponseVisitor) {
	visitor.Reply(s)
}

//it is not mandatory to register results in type system but it's certainly a good thing
func init() {
	inspectables.Register(stuffName, func() inspect.Inspectable { return new(stuff) })
}
