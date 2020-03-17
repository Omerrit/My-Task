package treeedit

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectable"
	"gerrit-share.lan/go/inspect/inspectables"
)

const (
	listFilesName      = packageName + ".ls"
	getStatsName       = packageName + ".getstats"
	saveName           = packageName + ".save"
	loadName           = packageName + ".load"
	getCurrentName     = packageName + ".get"
	editPositionName   = packageName + ".editpos"
	deletePositionName = packageName + ".delpos"
	getPositionName    = packageName + ".getpos"
	editDivisionName   = packageName + ".editdiv"
	deleteDivisionName = packageName + ".deldiv"
	newDivisionName    = packageName + ".newdiv"
	newPositionName    = packageName + ".newpos"
	getDivisionName    = packageName + ".getdiv"
)

type itemIdRequest struct {
	id string
}

func (d *itemIdRequest) Inspect(i *inspect.ObjectInspector) {
	i.String(&d.id, "path", true, "resource path")
}

type optionalItemIdRequest struct {
	id string
}

func (o *optionalItemIdRequest) Inspect(i *inspect.ObjectInspector) {
	i.String(&o.id, "path", false, "resource path")
}

type getCurrent struct {
	optionalItemIdRequest
	depth int
}

func (g *getCurrent) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("get_current", "get currently loaded tree")
	g.optionalItemIdRequest.Inspect(objectInspector)
	objectInspector.Int(&g.depth, "depth", false, "resulting tree depth")
	objectInspector.End()
}

type getStats struct {
	inspectable.DummyInspectable
}

type fileAction struct {
	file string
}

func (l *fileAction) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("", "load data from provided file")
	l.Embed(objectInspector)
	objectInspector.End()
}

func (l *fileAction) Embed(i *inspect.ObjectInspector) {
	i.String(&l.file, "f", true, "file name")
}

type loadFile struct {
	fileAction
}

type saveFile struct {
	optionalItemIdRequest
	fileAction
}

func (s *saveFile) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("", "save changes to provided file")
	s.optionalItemIdRequest.Inspect(objectInspector)
	s.Embed(objectInspector)
	objectInspector.End()
}

type positionDetails struct {
	isSuperior bool
	name       string
	position   string
}

func (p *positionDetails) Inspect(i *inspect.ObjectInspector) {
	i.Bool(&p.isSuperior, "boss", true, "this is a boss")
	i.String(&p.name, "name", false, "person name")
	i.String(&p.position, "pos", false, "position name")
}

type editPosition struct {
	optionalItemIdRequest
	positionDetails
	isEmpty bool
}

func (e *editPosition) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("", "edit provided position's data")
	e.optionalItemIdRequest.Inspect(objectInspector)
	objectInspector.Bool(&e.isEmpty, "empty", true, "position is empty")
	e.positionDetails.Inspect(objectInspector)
	objectInspector.End()
}

type divisionDetails struct {
	name string
}

func (d *divisionDetails) Inspect(i *inspect.ObjectInspector) {
	i.String(&d.name, "name", true, "division name")
}

type editDivision struct {
	optionalItemIdRequest
	divisionDetails
}

func (e *editDivision) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("", "edit provided division's data")
	e.optionalItemIdRequest.Inspect(objectInspector)
	e.divisionDetails.Inspect(objectInspector)
	objectInspector.End()
}

type newPosition struct {
	itemIdRequest
	positionDetails
}

func (n *newPosition) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("", "create new position with provided information inside provided division")
	n.itemIdRequest.Inspect(objectInspector)
	n.positionDetails.Inspect(objectInspector)
	objectInspector.End()
}

type newDivision struct {
	itemIdRequest
	divisionDetails
}

func (n *newDivision) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("", "create new position with provided information inside provided division")
	n.itemIdRequest.Inspect(objectInspector)
	n.divisionDetails.Inspect(objectInspector)
	objectInspector.End()
}

type deleteDivision struct {
	itemIdRequest
}

func (d *deleteDivision) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("", "delete provided division")
	d.itemIdRequest.Inspect(objectInspector)
	objectInspector.End()
}

type deletePosition struct {
	itemIdRequest
}

func (d *deletePosition) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("", "delete provided position")
	d.itemIdRequest.Inspect(objectInspector)
	objectInspector.End()
}

type divisionInfoRequest struct {
	optionalItemIdRequest
}

func (d *divisionInfoRequest) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("", "get information regarding provided division")
	d.optionalItemIdRequest.Inspect(objectInspector)
	objectInspector.End()
}

type positionInfoRequest struct {
	itemIdRequest
}

func (p *positionInfoRequest) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("", "get information regarding provided position")
	p.itemIdRequest.Inspect(objectInspector)
	objectInspector.End()
}

type listFiles struct {
	inspectable.DummyInspectable
}

func init() {
	inspectables.RegisterDescribed(listFilesName, func() inspect.Inspectable { return new(listFiles) },
		"get available file list")
	inspectables.RegisterDescribed(getStatsName, func() inspect.Inspectable { return new(getStats) },
		"get statistic of currently loaded tree")
	inspectables.Register(saveName, func() inspect.Inspectable { return new(saveFile) })
	inspectables.Register(loadName, func() inspect.Inspectable { return new(loadFile) })
	inspectables.Register(getCurrentName, func() inspect.Inspectable { return &getCurrent{depth: -1} })
	inspectables.Register(editPositionName, func() inspect.Inspectable { return new(editPosition) })
	inspectables.Register(deletePositionName, func() inspect.Inspectable { return new(deletePosition) })
	inspectables.Register(getPositionName, func() inspect.Inspectable { return new(positionInfoRequest) })
	inspectables.Register(editDivisionName, func() inspect.Inspectable { return new(editDivision) })
	inspectables.Register(deleteDivisionName, func() inspect.Inspectable { return new(deleteDivision) })
	inspectables.Register(newDivisionName, func() inspect.Inspectable { return new(newDivision) })
	inspectables.Register(newPositionName, func() inspect.Inspectable { return new(newPosition) })
	inspectables.Register(getDivisionName, func() inspect.Inspectable { return new(divisionInfoRequest) })
}
