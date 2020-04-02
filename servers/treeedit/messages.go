package treeedit

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

type itemIdRequest struct {
	id string
}

const itemIdRequestName = packageName + ".itemid"

func (d *itemIdRequest) Embed(i *inspect.ObjectInspector) {
	i.String(&d.id, "path", true, "resource path")
}

func (d *itemIdRequest) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(itemIdRequestName, "")
	{
		objectInspector.String(&d.id, "path", true, "resource path")
		objectInspector.End()
	}
}

func (d *itemIdRequest) Path() string {
	return d.id
}

type getCurrent struct {
	itemIdRequest
	depth int
}

const getCurrentName = packageName + ".get"

func (g *getCurrent) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(getCurrentName, "get currently loaded tree")
	{
		g.itemIdRequest.Embed(objectInspector)
		objectInspector.Int(&g.depth, "depth", false, "resulting tree depth")
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(getCurrentName, func() inspect.Inspectable { return &getCurrent{depth: -1} })
}

type getStats struct {
	inspect.EmptyObject
}

func init() {
	inspectables.RegisterDescribed(packageName+".getstats", func() inspect.Inspectable { return new(getStats) },
		"get statistics of currently loaded tree")
}

type fileAction struct {
	file string
}

const fileActionName = packageName + ".file"

func (l *fileAction) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(fileActionName, "file action")
	{
		l.Embed(objectInspector)
		objectInspector.End()
	}
}

func (l *fileAction) Embed(i *inspect.ObjectInspector) {
	i.String(&l.file, "f", true, "file name")
}

type loadFile struct {
	fileAction
}

func init() {
	inspectables.RegisterDescribed(packageName+".load", func() inspect.Inspectable { return new(loadFile) },
		"load data from provided file")
}

type saveFile struct {
	itemIdRequest
	fileAction
}

const saveFileName = packageName + ".save"

func (s *saveFile) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(saveFileName, "save changes to provided file")
	{
		s.itemIdRequest.Embed(objectInspector)
		s.fileAction.Embed(objectInspector)
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(saveFileName, func() inspect.Inspectable { return new(saveFile) })
}

type positionDetails struct {
	isSuperior bool
	name       string
	position   string
}

func (p *positionDetails) Embed(i *inspect.ObjectInspector) {
	i.Bool(&p.isSuperior, "boss", true, "this is a boss")
	i.String(&p.name, "name", false, "person name")
	i.String(&p.position, "pos", false, "position name")
}

type editPosition struct {
	itemIdRequest
	positionDetails
	isEmpty bool
}

const editPositionName = packageName + ".editpos"

func (e *editPosition) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(editPositionName, "edit provided position's data")
	{
		e.itemIdRequest.Embed(objectInspector)
		objectInspector.Bool(&e.isEmpty, "empty", true, "position is empty")
		e.positionDetails.Embed(objectInspector)
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(editPositionName, func() inspect.Inspectable { return new(editPosition) })
}

type divisionDetails struct {
	name string
}

func (d *divisionDetails) Embed(i *inspect.ObjectInspector) {
	i.String(&d.name, "name", true, "division name")
}

type editDivision struct {
	itemIdRequest
	divisionDetails
}

const editDivisionName = packageName + ".editdiv"

func (e *editDivision) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(editDivisionName, "edit provided division's data")
	{
		e.itemIdRequest.Embed(objectInspector)
		e.divisionDetails.Embed(objectInspector)
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(editDivisionName, func() inspect.Inspectable { return new(editDivision) })
}

type newPosition struct {
	itemIdRequest
	positionDetails
}

const newPositionName = packageName + ".newpos"

func (n *newPosition) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(newPositionName, "create new position with provided information inside provided division")
	{
		n.itemIdRequest.Embed(objectInspector)
		n.positionDetails.Embed(objectInspector)
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(newPositionName, func() inspect.Inspectable { return new(newPosition) })
}

type newDivision struct {
	itemIdRequest
	divisionDetails
}

const newDivisionName = packageName + ".newdiv"

func (n *newDivision) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(newDivisionName, "create new division with provided information inside provided division")
	{
		n.itemIdRequest.Embed(objectInspector)
		n.divisionDetails.Embed(objectInspector)
		objectInspector.End()
	}
}

func init() {
	inspectables.Register(newDivisionName, func() inspect.Inspectable { return new(newDivision) })
}

type deleteDivision struct {
	itemIdRequest
}

func init() {
	inspectables.RegisterDescribed(packageName+".deldiv", func() inspect.Inspectable { return new(deleteDivision) },
		"delete provided division")
}

type deletePosition struct {
	itemIdRequest
}

func init() {
	inspectables.RegisterDescribed(packageName+".delpos", func() inspect.Inspectable { return new(deletePosition) },
		"delete provided position")
}

type divisionInfoRequest struct {
	itemIdRequest
}

func init() {
	inspectables.RegisterDescribed(packageName+".getdiv", func() inspect.Inspectable { return new(divisionInfoRequest) },
		"get information regarding provided division")
}

type positionInfoRequest struct {
	itemIdRequest
}

func init() {
	inspectables.RegisterDescribed(packageName+".getpos", func() inspect.Inspectable { return new(positionInfoRequest) },
		"get information regarding provided position")
}

type listFiles struct {
	inspect.EmptyObject
}

func init() {
	inspectables.RegisterDescribed(packageName+".ls", func() inspect.Inspectable { return new(listFiles) },
		"get available file list")
}
