package treeedit

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectable/arrays"
	"gerrit-share.lan/go/inspect/inspectable/maps"
	"gerrit-share.lan/go/inspect/inspectables"
	"gerrit-share.lan/go/inspect/json/fromjson"
	"gerrit-share.lan/go/utils/sets"
	"io/ioutil"
	"sort"
	"time"
)

const (
	positionName          = packageName + ".pos"
	divisionName          = packageName + ".div"
	genericDivisionName   = packageName + ".gendiv"
	divisionShortenedName = packageName + ".divshort"
	divisionArrayName     = packageName + ".divarray"
)

type Position struct {
	Id         string
	IsSuperior bool
	IsEmpty    bool
	Name       string
	Position   string
	StartDate  string
	DbId       int64
	Info       maps.MapStringString
}

const TimeFormatString = "2 Jan 2006 15:04"

func NewPosition(isSuperior bool, name string, positionName string) *Position {
	return &Position{"", isSuperior, len(name) == 0, name, positionName, time.Now().Format(TimeFormatString), 0, make(maps.MapStringString)}
}

type Division struct {
	Id           string
	Name         string
	Companies    arrays.StringArray
	Positions    Positions
	Divisions    Divisions
	DbId         int64
	StartDate    string
	Info         maps.MapStringString
	nextPosition int64
	nextDivision int64
}

func NewDivision(name string) *Division {
	return &Division{"", name, nil, make(Positions), make(Divisions), 0, time.Now().Format(TimeFormatString), make(maps.MapStringString), 0, 0}
}

func (d *Division) AddDivision(div *Division) {
	name := EncodeIntId(d.nextDivision)
	d.nextDivision++
	if len(d.Id) > 0 {
		div.Id = d.Id + "." + name
	} else {
		div.Id = name
	}
	d.Divisions[name] = div
}

func (d *Division) AddPosition(pos *Position) {
	name := EncodeIntId(d.nextPosition)
	d.nextPosition++
	if len(d.Id) > 0 {
		pos.Id = d.Id + "." + name
	} else {
		pos.Id = name
	}
	d.Positions[name] = pos
}

func (d *Division) RemoveDivision(div *Division) {
	if d == nil {
		return
	}
	for name, division := range d.Divisions {
		if division.Id == div.Id {
			delete(d.Divisions, name)
			return
		}
	}
}

func (d *Division) RemoveDivisionByShortName(name string) {
	if d == nil {
		return
	}
	delete(d.Divisions, name)
}

func (d *Division) RemovePosition(pos *Position) {
	if d == nil {
		return
	}
	for name, position := range d.Positions {
		if position.Id == pos.Id {
			delete(d.Positions, name)
			return
		}
	}
}

func (d *Division) RemovePositionByShortName(name string) {
	if d == nil {
		return
	}
	delete(d.Positions, name)
}

func (d *Division) FixIds() {
	if d == nil {
		return
	}
	prefix := ""
	if len(d.Id) > 0 {
		prefix = d.Id + "."
	}
	for name, pos := range d.Positions {
		pos.Id = prefix + name
	}
	for name, div := range d.Divisions {
		div.Id = prefix + name
		div.FixIds()
	}
}

func (p *Position) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(positionName, "position object")
	objectInspector.String(&p.Id, "id", true, "position id")
	objectInspector.Bool(&p.IsSuperior, "is_superior", true, "division boss")
	objectInspector.Bool(&p.IsEmpty, "is_empty", true, "empty position")
	objectInspector.String(&p.Name, "name", true, "person name")
	objectInspector.String(&p.Position, "position", true, "position name")
	objectInspector.String(&p.StartDate, "start_date", false, "start date")
	objectInspector.Int64(&p.DbId, "db_id", false, "db id")
	genericInspector := objectInspector.Value("info", false, "additional information")
	if genericInspector != nil {
		p.Info.Inspect(genericInspector)
	}
	objectInspector.End()
}

func (p *Position) Visit(v actors.ResponseVisitor) {
	v.Reply(p)
}

func (d *Division) CutTree(depth int) *Division {
	if d == nil {
		return nil
	}
	node := NewDivision(d.Name)
	node.Id = d.Id
	node.Companies = d.Companies
	node.DbId = d.DbId
	node.Info = d.Info
	if depth != 0 {
		node.Positions = d.Positions
		node.Divisions = make(Divisions, len(d.Divisions))
		for key, v := range d.Divisions {
			node.Divisions[key] = v.CutTree(depth - 1)
		}
	}
	return node
}

func (d *Division) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(divisionName, "division object")
	objectInspector.String(&d.Id, "id", true, "division id")
	objectInspector.String(&d.Name, "name", true, "division name")
	d.Companies.Inspect(objectInspector.Value("companies", true, "division companies"))
	d.Positions.Inspect(objectInspector.Value("positions", true, "underlying positions list"))
	d.Divisions.Inspect(objectInspector.Value("divisions", true, "underlying divisions list"))
	objectInspector.Int64(&d.nextPosition, "next_position", true, "")
	objectInspector.Int64(&d.nextDivision, "next_division", true, "")
	objectInspector.Int64(&d.DbId, "db_id", false, "db id")
	objectInspector.String(&d.StartDate, "start_date", false, "start date")
	genericInspector := objectInspector.Value("info", false, "additional information")
	if genericInspector != nil {
		d.Info.Inspect(genericInspector)
	}
	objectInspector.End()
}

func (d *Division) Visit(v actors.ResponseVisitor) {
	v.Reply(d)
}

type genericDivision Division

func (d *genericDivision) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(genericDivisionName, "division object")
	objectInspector.String(&d.Id, "id", true, "division id")
	objectInspector.String(&d.Name, "name", true, "division name")
	objectInspector.Int64(&d.DbId, "db_id", false, "db id")
	objectInspector.String(&d.StartDate, "start_date", false, "start date")
	genericInspector := objectInspector.Value("info", false, "additional information")
	if genericInspector != nil {
		d.Info.Inspect(genericInspector)
	}
	d.Companies.Inspect(objectInspector.Value("companies", true, "division companies"))
	divisionAsArray(*d).Inspect(objectInspector.Value("children", true, "underlying divisions and positions list"))
	objectInspector.End()
}

func (d *genericDivision) Visit(v actors.ResponseVisitor) {
	v.Reply(d)
}

type divisionAsArray Division

func (g divisionAsArray) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(divisionArrayName, genericDivisionName, "array of divisions")
	if arrayInspector.IsReading() {
		return
	}
	arrayInspector.SetLength(len(g.Divisions) + len(g.Positions))
	for _, v := range g.Divisions {
		(*genericDivision)(v).Inspect(arrayInspector.Value())
	}
	for _, v := range g.Positions {
		v.Inspect(arrayInspector.Value())
	}
	arrayInspector.End()
}

type DivisionShortened Division

func (d *DivisionShortened) Inspect(i *inspect.GenericInspector) {
	if !i.IsReading() {
		positionsIds := make(arrays.StringArray, len(d.Positions))
		var counter int
		for name := range d.Positions {
			positionsIds[counter] = name
			counter++
		}
		counter = 0
		divisionsIds := make(arrays.StringArray, len(d.Divisions))
		for name := range d.Divisions {
			divisionsIds[counter] = name
			counter++
		}
		objectInspector := i.Object(divisionShortenedName, "division object")
		objectInspector.String(&d.Id, "id", true, "division id")
		objectInspector.String(&d.Name, "name", true, "division name")
		d.Companies.Inspect(objectInspector.Value("companies", true, "division companies"))
		positionsIds.Inspect(objectInspector.Value("positions", true, "underlying positions ids"))
		divisionsIds.Inspect(objectInspector.Value("divisions", true, "underlying divisions ids"))
		objectInspector.End()
	}
}

func LoadFile(name string) (*Division, error) {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var top OldStructTop
	reader := inspect.NewGenericInspector(fromjson.NewInspector(content, 0))
	top.Inspect(reader)
	if len(top.List) == 0 {
		var div Division
		reader := inspect.NewGenericInspector(fromjson.NewInspector(content, 0))
		div.Inspect(reader)
		if reader.GetError() != nil {
			var div1 Division1
			reader := inspect.NewGenericInspector(fromjson.NewInspector(content, 0))
			div1.Inspect(reader)
			if reader.GetError() != nil {
				return nil, reader.GetError()
			}
			return div1.Upgrade(), nil
		}
		return &div, nil
	}
	result := NewDivision("")
	var companies sets.String
	for _, item := range top.List {
		div := ConvertDivision(item)
		div.Name = item.Details.Rdc
		result.AddDivision(div)
		for _, company := range div.Companies {
			companies.Add(company)
		}
	}
	for company := range companies {
		result.Companies = append(result.Companies, company)
	}
	sort.Strings(result.Companies)
	return result, nil
}

func init() {
	inspectables.Register(positionName, func() inspect.Inspectable { return new(Position) })
	inspectables.Register(divisionName, func() inspect.Inspectable { return new(DivisionShortened) })
	inspectables.Register(genericDivisionName, func() inspect.Inspectable { return new(genericDivision) })
	inspectables.Register(divisionShortenedName, func() inspect.Inspectable { return new(DivisionShortened) })
	inspectables.Register(divisionArrayName, func() inspect.Inspectable { return new(divisionAsArray) })
}
