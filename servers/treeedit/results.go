package treeedit

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

type Id string

const idName = packageName + ".id"

func (i Id) Visit(v actors.ResponseVisitor) {
	v.ReplyString((string)(i))
}

func (i *Id) Inspect(inspector *inspect.GenericInspector) {
	inspector.String((*string)(i), idName, "id")
}

func init() {
	inspectables.Register(idName, func() inspect.Inspectable { return new(Id) })
}

type Statistics struct {
	Divisions int
	Positions int
	Companies int
}

const StatisticsName = packageName + ".stats"

func (s *Statistics) Visit(v actors.ResponseVisitor) {
	v.Reply(s)
}

func (s *Statistics) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(StatisticsName, "divisions, positions and companies statistic")
	{
		objectInspector.Int(&s.Divisions, "divisions", false, "number of underlying divisions")
		objectInspector.Int(&s.Positions, "positions", false, "number of underlying positions")
		objectInspector.Int(&s.Companies, "companies", false, "number of underlying companies")
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(StatisticsName, func() inspect.Inspectable { return new(Statistics) })
}

type files []string

const filesName = packageName + ".files"

func (f files) Visit(v actors.ResponseVisitor) {
	v.Reply(f)
}

func (f files) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(filesName, "string", "list of available files")
	{
		if arrayInspector.IsReading() {
			return
		}
		arrayInspector.SetLength(len(f))
		for _, v := range f {
			arrayInspector.String(&v)
		}
		arrayInspector.End()
	}
}

func (f *files) Add(name string) {
	*f = append(*f, name)
}

func init() {
	inspectables.Register(filesName, func() inspect.Inspectable { return files{} })
}
