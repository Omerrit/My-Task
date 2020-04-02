package actors

import (
	"gerrit-share.lan/go/inspect"
)

type PortableType struct {
	typeName string
	sample   inspect.Inspectable
}

const PortableTypeName = packageName + ".ptype"

func (p PortableType) TypeName() string {
	return p.typeName
}

func (p PortableType) Sample() inspect.Inspectable {
	return p.sample
}

func (p *PortableType) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(PortableTypeName, "serializable value with type name to be read")
	{
		o.String(&p.typeName, "name", true, "type name")
		if o.IsReading() {
			p.sample = o.GetExtraValue().(TypeSystem).Creator(p.typeName)()
		}
		p.sample.Inspect(o.Value("sample", true, "payload"))
		o.End()
	}
}

type PortableValues []PortableType

const PortableValuesName = packageName + ".pvalues"

func (p *PortableValues) Inspect(inspector *inspect.GenericInspector) {
	a := inspector.Array(PortableValuesName, PortableTypeName, "")
	{
		if a.IsReading() {
			length := a.GetLength()
			if length > cap(*p) {
				*p = make(PortableValues, length)
			} else {
				*p = (*p)[:length]
			}
		}
		for i := range *p {
			(*p)[i].Inspect(a.Value())
		}
		a.End()
	}
}
