package example

//commands go here if there's a lot of them (more than a couple or they just have lots of methods)

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

type pathHolder struct {
	path string
}

const pathHolderName = packageName + ".path"

//there is a special 'path' property that gets mapped to REST API URL path when using http server
//this function should be implemented and 'path' property be declared
func (p *pathHolder) Path() string {
	return p.path
}

func (p *pathHolder) Embed(inspector *inspect.ObjectInspector) {
	//declare 'path' property
	inspector.String(&p.path, "path", true, "command path")
}

func (p *pathHolder) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(pathHolderName, "command path container")
	{
		p.Embed(o)
		o.End()
	}
}

//we don't plan to use pathHolder as command,message or result so we don't register it

type getStuff struct {
	pathHolder
}

//we don't export getStuff type, it has no methods and we don't plan to make containers with it
//so we don't make type name constant and use RegisterDescribed to register this command in the type system
func init() {
	inspectables.RegisterDescribed(packageName+".get", func() inspect.Inspectable { return new(getStuff) },
		"retrieves some stuff")
}
