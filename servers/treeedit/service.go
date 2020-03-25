package treeedit

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/published"
	"gerrit-share.lan/go/actors/starter"
	"gerrit-share.lan/go/debug"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectable/arrays"
	"gerrit-share.lan/go/inspect/json/tojson"
	"gerrit-share.lan/go/utils/flags"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

//API:
// - load file (convert old one if necessary)
// - get current json
// - save to file
// - set parameters by id
// - add person, get id
// - add division, get id and superior id

const (
	serviceName = "tree_edit"
)

type divisionInfo struct {
	*Division
	parent *Division
}

func (d divisionInfo) Visit(v actors.ResponseVisitor) {
	v.Reply((*DivisionShortened)(d.Division))
}

type positionInfo struct {
	*Position
	parent *Division
}

type TreeEditor struct {
	actors.Actor
	root      *Division
	workDir   string
	divisions map[string]divisionInfo
	positions map[string]positionInfo
	name      string
}

func (t *TreeEditor) MakeBehaviour() actors.Behaviour {
	log.Println(t.name, "started")
	divisionSample := &genericDivision{Companies: arrays.StringArray{}, Positions: map[string]*Position{"": &Position{}}, Divisions: map[string]*Division{"": &Division{}}}
	var b actors.Behaviour
	b.Name = serviceName
	var handle starter.Handle
	handle.Acquire(t, handle.DependOn, t.Quit)
	b.AddCommand(new(listFiles), t.listFiles).Result(files{""})
	b.AddCommand(new(getStats), t.getStats).Result(new(Statistics))
	b.AddCommand(new(saveFile), t.save)
	b.AddCommand(new(loadFile), t.load).Result(divisionSample)
	b.AddCommand(&getCurrent{depth: -1}, t.getCurrentNode).Result(divisionSample)
	b.AddCommand(new(editPosition), t.editPosition)
	b.AddCommand(new(deletePosition), t.deletePosition)
	b.AddCommand(new(positionInfoRequest), t.getPosition).Result(new(Position))
	b.AddCommand(new(editDivision), t.editDivision)
	b.AddCommand(new(deleteDivision), t.deleteDivision)
	b.AddCommand(new(newDivision), t.newDivision).Result(new(Id))
	b.AddCommand(new(newPosition), t.newPosition).Result(new(Id))
	b.AddCommand(new(divisionInfoRequest), t.getDivision).Result((*DivisionShortened)(divisionSample))
	t.SetPanicProcessor(func(err errors.StackTraceError) {
		log.Println("panic:", err, err.StackTrace())
		t.Quit(err)
	})
	published.Publish(t, t.Quit)
	return b
}

func (t *TreeEditor) Shutdown() error {
	log.Println(t.name, "shut down")
	return nil
}

func (t *TreeEditor) regenerateDatabase() {
	t.divisions = make(map[string]divisionInfo)
	t.positions = make(map[string]positionInfo)
	if t.root != nil {
		t.root.FixIds()
		t.divisions[t.root.Id] = divisionInfo{t.root, nil}
		t.indexDivision(t.root)
	} else {
		t.root = NewDivision("")
		t.divisions[t.root.Id] = divisionInfo{t.root, nil}
	}
}

func (t *TreeEditor) indexDivision(div *Division) {
	if div == nil {
		return
	}
	for _, pos := range div.Positions {
		if pos != nil {
			t.positions[pos.Id] = positionInfo{pos, div}
		}
	}
	for _, nextDiv := range div.Divisions {
		t.divisions[nextDiv.Id] = divisionInfo{nextDiv, div}
		t.indexDivision(nextDiv)
	}
}

func (t *TreeEditor) listFiles(_ interface{}) (actors.Response, error) {
	infos, err := ioutil.ReadDir(t.workDir)
	if err != nil {
		return nil, err
	}
	var result files
	for _, info := range infos {
		if !info.IsDir() {
			result.Add(info.Name())
		}
	}
	return result, nil
}

func (t *TreeEditor) getStats(_ interface{}) (actors.Response, error) {
	stat := &Statistics{
		Positions: len(t.positions),
		Divisions: len(t.divisions)}
	if t.root != nil {
		stat.Companies = len(t.root.Companies)
	}
	return stat, nil
}

func (t *TreeEditor) save(cmd interface{}) (actors.Response, error) {
	saveCmd := cmd.(*saveFile)
	div, err := t.getCurrent(saveCmd.id)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(path.Join(t.workDir, path.Base(saveCmd.file)))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	inspector := &tojson.Inspector{}
	serializer := inspect.NewGenericInspector(inspector)
	div.Inspect(serializer)
	if serializer.GetError() != nil {
		return nil, serializer.GetError()
	}
	_, err = file.Write(inspector.Output())
	return nil, err
}

func (t *TreeEditor) getCurrent(id string) (*Division, error) {
	div := t.divisions[id]
	if div.Division == nil {
		pos := t.positions[id]
		if pos.Position == nil {
			return nil, ErrInvalidId
		}
		div.Division = pos.parent
	}
	return div.Division, nil
}

func (t *TreeEditor) load(cmd interface{}) (actors.Response, error) {
	loadCmd := cmd.(*loadFile)
	var err error
	t.root, err = LoadFile(path.Join(t.workDir, path.Base(loadCmd.file)))
	if err != nil {
		return nil, err
	}
	t.regenerateDatabase()
	return (*genericDivision)(t.root.CutTree(0)), nil
}

func (t *TreeEditor) getCurrentNode(cmd interface{}) (actors.Response, error) {
	getCmd := cmd.(*getCurrent)
	div := t.findDivisionByPath(getCmd.id)
	if div == nil {
		return nil, ErrInvalidId
	}
	return (*genericDivision)(div.CutTree(getCmd.depth)), nil
}

func (t *TreeEditor) findDivisionByPath(path string) *Division {
	if len(path) == 0 {
		return t.root
	}
	div := t.root
	if div == nil {
		return nil
	}
	for {
		next := strings.IndexByte(path, '.')
		if next != -1 {
			div = div.Divisions[path[:next]]
			if div == nil {
				return nil
			}
			path = path[(next + 1):]
		} else {
			return div.Divisions[path]
		}
	}
}

func (t *TreeEditor) editPosition(cmd interface{}) (actors.Response, error) {
	editCmd := cmd.(*editPosition)
	pos := t.positions[editCmd.id]
	if pos.Position == nil {
		return nil, ErrInvalidPosition
	}
	pos.IsEmpty = editCmd.isEmpty
	pos.IsSuperior = editCmd.isSuperior
	if len(editCmd.name) > 0 {
		pos.Name = editCmd.name
	}
	if len(editCmd.position) > 0 {
		pos.Position.Position = editCmd.position
	}
	return nil, nil
}

func (t *TreeEditor) deletePosition(cmd interface{}) (actors.Response, error) {
	deleteCmd := cmd.(*deletePosition)
	pos := t.positions[deleteCmd.id]
	if pos.Position == nil {
		return nil, ErrInvalidPosition
	}
	if pos.parent == nil {
		return nil, ErrPositionIsRoot
	}
	pos.parent.RemovePosition(pos.Position)
	delete(t.positions, pos.Id)
	return nil, nil
}

func (t *TreeEditor) getPosition(cmd interface{}) (actors.Response, error) {
	getCmd := cmd.(*positionInfoRequest)
	pos, ok := t.positions[getCmd.id]
	if !ok {
		return nil, ErrInvalidPosition
	}
	return pos, nil
}

func (t *TreeEditor) editDivision(cmd interface{}) (actors.Response, error) {
	editCmd := cmd.(*editDivision)
	div := t.divisions[editCmd.id]
	if div.Division == nil {
		return nil, ErrInvalidDivision
	}
	div.Name = editCmd.name
	return nil, nil
}

func (t *TreeEditor) deleteDivision(cmd interface{}) (actors.Response, error) {
	deleteCmd := cmd.(*deleteDivision)
	div := t.divisions[deleteCmd.id]
	if div.Division == nil {
		return nil, ErrInvalidDivision
	}
	if div.parent == nil {
		return nil, ErrCannotRemoveRoot
	}
	t.removeDivisionFromIndex(div.Division)
	div.parent.RemoveDivision(div.Division)
	delete(t.divisions, div.Id)
	return nil, nil
}

func (t *TreeEditor) removeDivisionFromIndex(div *Division) {
	for _, pos := range div.Positions {
		delete(t.positions, pos.Id)
	}
	for _, nextDiv := range div.Divisions {
		t.removeDivisionFromIndex(nextDiv)
		delete(t.divisions, nextDiv.Id)
	}
}

func (t *TreeEditor) newDivision(cmd interface{}) (actors.Response, error) {
	newCmd := cmd.(*newDivision)
	div := t.divisions[newCmd.id]
	if div.Division == nil {
		return nil, ErrInvalidParent
	}
	newDiv := NewDivision(newCmd.name)
	div.AddDivision(newDiv)
	t.divisions[newDiv.Id] = divisionInfo{newDiv, div.Division}
	return Id(newDiv.Id), nil
}

func (t *TreeEditor) newPosition(cmd interface{}) (actors.Response, error) {
	newCmd := cmd.(*newPosition)
	div := t.divisions[newCmd.id]
	if div.Division == nil {
		return nil, ErrInvalidParent
	}
	pos := NewPosition(newCmd.isSuperior, newCmd.name, newCmd.position)
	div.AddPosition(pos)
	t.positions[pos.Id] = positionInfo{pos, div.Division}
	return Id(pos.Id), nil
}

func (t *TreeEditor) getDivision(cmd interface{}) (actors.Response, error) {
	getCmd := cmd.(*divisionInfoRequest)
	div := t.divisions[getCmd.id]
	if div.Division == nil {
		return nil, ErrInvalidDivision
	}
	return t.divisions[getCmd.id], nil
}

func init() {
	const name = "tree_editor"
	var preloadFile string
	dir, _ := os.Getwd()
	starter.SetCreator(name, func(actor *actors.Actor, name string) (actors.ActorService, error) {
		ed := &TreeEditor{workDir: dir, name: name}
		if len(preloadFile) != 0 {
			var err error
			ed.root, err = LoadFile(preloadFile)
			if err != nil {
				debug.Printf("Error opening file: %#v, %v", err, err)
			}
		}
		ed.DependOn(actor.Service())
		return actor.System().Spawn(ed), nil
	})
	starter.SetFlagInitializer(name, func() {
		flags.StringFlag(&dir, "dir", "work directory")
		flags.StringFlag(&preloadFile, "file", "load this file")
	})
}
