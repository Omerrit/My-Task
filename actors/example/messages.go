package example

//this file contains all message types that service accepts
//it also may contain command types

//all message and command types must be inspectable and must be registered in type system (inspectables)
//durung package init

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

//packageName string should be declared
const packageName = "example"

type embeddable struct {
	name string
}

//implement Embed if you want this type to be easily embedded (be presented as part of) into other types
//if you want this type to also be inspectable then inspect should use embed
func (e *embeddable) Embed(inspector *inspect.ObjectInspector) {
	inspector.String(&e.name, "name", true, "name field")
}

type message struct {
	name embeddable
}

//type name constant should be present for all inspectable types that have their own Inspect method
//and for all exported inspectable types
//type name constant should immediately follow type definition and named typeName
//if type is exported then its name constant should be exported too

//use this constant when implementing inspect for this type
//and when registering this type

//type name should be unique but preferably short
//use snake_case for naming, it's permissible to shorten words and omit underscore
const messageName = packageName + ".message"

func (m *message) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(messageName, "example message type")
	//use braces to mark code where you use newly obtained type inspector
	{
		m.name.Embed(o)
		//always call End() for a type inspector
		//this should be the last sentence of a block
		o.End()
	}
}

//there should be one init() per type and it should be placed after all type methods
//one init that register several types is possible if types are short
//(up to 30 lines combined including type names and extra empty lines, preferably 20-25 lines),
//have no methods and are declared close to each other
func init() {
	//use Register to register inspectable objects that have their own inspect()
	//and RegisterDescribed to register inspectable objects without their own inspect()
	inspectables.Register(messageName, func() inspect.Inspectable { return new(message) })
}
