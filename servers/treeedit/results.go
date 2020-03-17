package treeedit

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

const (
	statisticsName = packageName + ".stats"
	filesName      = packageName + ".files"
	idName         = packageName + ".id"
)

type Id string

func (i Id) Visit(v actors.ResponseVisitor) {
	v.ReplyString((string)(i))
}

func (i *Id) Inspect(inspector *inspect.GenericInspector) {
	inspector.String((*string)(i), idName, "id")
}

type Statistics struct {
	Divisions int
	Positions int
	Companies int
}

func (s *Statistics) Visit(v actors.ResponseVisitor) {
	v.Reply(s)
}

func (s *Statistics) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(statisticsName, "divisions, positions and companies statistic")
	objectInspector.Int(&s.Divisions, "divisions", false, "number of underlying divisions")
	objectInspector.Int(&s.Positions, "positions", false, "number of underlying positions")
	objectInspector.Int(&s.Companies, "companies", false, "number of underlying companies")
	objectInspector.End()
}

type files []string

func (f files) Visit(v actors.ResponseVisitor) {
	v.Reply(f)
}

func (f files) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(filesName, "", "list of available files")
	if arrayInspector.IsReading() {
		return
	}
	arrayInspector.SetLength(len(f))
	for _, v := range f {
		arrayInspector.String(&v)
	}
	arrayInspector.End()
}

func (f *files) Add(name string) {
	*f = append(*f, name)
}

func init() {
	inspectables.Register(filesName, func() inspect.Inspectable { return files{} })
	inspectables.Register(statisticsName, func() inspect.Inspectable { return new(Statistics) })
	inspectables.Register(idName, func() inspect.Inspectable { return new(Id) })
}
