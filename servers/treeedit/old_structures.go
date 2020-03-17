package treeedit

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectable/arrays"
	"gerrit-share.lan/go/utils/sets"
	"sort"
	"strings"
)

type OldStructPersons []OldStructPerson

func (o *OldStructPersons) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array("list", "", "list")
	if !arrayInspector.IsReading() {
		arrayInspector.SetLength(len(*o))
	} else {
		*o = make(OldStructPersons, arrayInspector.GetLength())
	}
	for index := range *o {
		(*o)[index].Inspect(arrayInspector.Value())
	}
	arrayInspector.End()
}

type OldStructTop struct {
	List OldStructPersons
}

func (o *OldStructTop) Inspect(i *inspect.GenericInspector) {
	o.List.Inspect(i)
}

type OldStructPerson struct {
	Key      string
	Id       int
	Name     string
	Position string
	Children OldStructPersons
	Details  OldStructDetail
}

func (o *OldStructPerson) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("person", "person")
	objectInspector.String(&o.Key, "key", true, "key")
	objectInspector.Int(&o.Id, "id", true, "id")
	objectInspector.String(&o.Name, "fio", true, "fio")
	objectInspector.String(&o.Position, "position", true, "position")
	o.Children.Inspect(objectInspector.Value("children", true, "children"))
	o.Details.Inspect(objectInspector.Value("detail", true, "detail"))
	objectInspector.End()
}

type OldStructDetail struct {
	CompanyCode string
	Company     string
	Rdc         string
	ChannelName string
	ProjectName string
}

func (o *OldStructDetail) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("old_detail", "old detail")
	objectInspector.String(&o.CompanyCode, "companyCode", true, "company code")
	objectInspector.String(&o.Company, "company", true, "company")
	objectInspector.String(&o.Rdc, "rdc", true, "rdc")
	objectInspector.String(&o.ChannelName, "channelName", true, "channel name")
	objectInspector.String(&o.ProjectName, "projectName", true, "project name")
	objectInspector.End()
}

type Division1 struct {
	Id        string
	Name      string
	Companies arrays.StringArray
	Positions positionArray
	Divisions division1Array
}

func (d *Division1) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object("division", "division object with uid")
	objectInspector.String(&d.Id, "id", true, "division id")
	objectInspector.String(&d.Name, "name", true, "division name")
	d.Companies.Inspect(objectInspector.Value("companies", true, "division companies"))
	d.Positions.Inspect(objectInspector.Value("positions", true, "underlying positions list"))
	d.Divisions.Inspect(objectInspector.Value("divisions", true, "underlying divisions list"))
	objectInspector.End()
}

func (d *Division1) Upgrade() *Division {
	div := NewDivision(d.Name)
	div.Companies = d.Companies
	for _, pos := range d.Positions {
		div.AddPosition(&pos)
	}
	for _, div1 := range d.Divisions {
		div.AddDivision(div1.Upgrade())
	}
	div.FixIds()
	return div
}

func ConvertPerson(person OldStructPerson) *Position {
	pos := NewPosition(len(person.Children) > 0, person.Name, person.Position)
	pos.IsEmpty = strings.Contains(person.Name, "Вакансия")
	return pos
}

func ConvertDivision(person OldStructPerson) *Division {
	div := NewDivision("")
	div.AddPosition(ConvertPerson(person))
	var companies sets.String
	for _, child := range person.Children {
		if len(child.Children) == 0 {
			div.AddPosition(ConvertPerson(child))
			companies.Add(child.Details.Company)
		} else {
			division := ConvertDivision(child)
			for _, comp := range division.Companies {
				companies.Add(comp)
			}
			div.AddDivision(division)
		}
	}
	for comp := range companies {
		div.Companies = append(div.Companies, comp)
	}
	sort.Strings(div.Companies)
	return div
}
