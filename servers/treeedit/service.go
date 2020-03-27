package treeedit

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/log"
	"gerrit-share.lan/go/actors/plugins/published"
	"gerrit-share.lan/go/actors/starter"
	"gerrit-share.lan/go/debug"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectable/arrays"
	"gerrit-share.lan/go/inspect/json/tojson"
	"gerrit-share.lan/go/utils/flags"
	"io/ioutil"
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
	log.Logger
	root      *Division
	workDir   string
	divisions map[string]divisionInfo
	positions map[string]positionInfo
	name      string
}

func (t *TreeEditor) MakeBehaviour() actors.Behaviour {
	t.SetLogSource(t.name)
	t.Infoln(t.name, "started")
	divisionSample := &genericDivision{Companies: arrays.StringArray{}, Positions: map[string]*Position{"": &Position{}}, Divisions: map[string]*Division{"": &Division{}}}
	var behaviour actors.Behaviour
	behaviour.Name = serviceName
	var starterHandle starter.Handle
	starterHandle.Acquire(t, starterHandle.DependOn, t.Quit)
	behaviour.AddCommand(new(listFiles), func(cmd interface{}) (actors.Response, error) {
		return t.listFiles()
	}).Result(files{""})
	behaviour.AddCommand(new(getStats), func(cmd interface{}) (actors.Response, error) {
		return t.getStats(), nil
	}).Result(new(Statistics))
	behaviour.AddCommand(new(saveFile), func(cmd interface{}) (actors.Response, error) {
		return nil, t.save(cmd.(*saveFile))
	})
	behaviour.AddCommand(new(loadFile), func(cmd interface{}) (actors.Response, error) {
		return t.load(cmd.(*loadFile))
	}).Result(divisionSample)
	behaviour.AddCommand(&getCurrent{depth: -1}, func(cmd interface{}) (actors.Response, error) {
		return t.getCurrentNode(cmd.(*getCurrent))
	}).Result(divisionSample)
	behaviour.AddCommand(new(editPosition), func(cmd interface{}) (actors.Response, error) {
		return nil, t.editPosition(cmd.(*editPosition))
	})
	behaviour.AddCommand(new(deletePosition), func(cmd interface{}) (actors.Response, error) {
		return nil, t.deletePosition(cmd.(*deletePosition))
	})
	behaviour.AddCommand(new(positionInfoRequest), func(cmd interface{}) (actors.Response, error) {
		return t.getPosition(cmd.(*positionInfoRequest))
	}).Result(new(Position))
	behaviour.AddCommand(new(editDivision), func(cmd interface{}) (actors.Response, error) {
		return nil, t.editDivision(cmd.(*editDivision))
	})
	behaviour.AddCommand(new(deleteDivision), func(cmd interface{}) (actors.Response, error) {
		return nil, t.deleteDivision(cmd.(*deleteDivision))
	})
	behaviour.AddCommand(new(newDivision), func(cmd interface{}) (actors.Response, error) {
		return t.newDivision(cmd.(*newDivision))
	}).Result(new(Id))
	behaviour.AddCommand(new(newPosition), func(cmd interface{}) (actors.Response, error) {
		return t.newPosition(cmd.(*newPosition))
	}).Result(new(Id))
	behaviour.AddCommand(new(divisionInfoRequest), func(cmd interface{}) (actors.Response, error) {
		return t.getDivision(cmd.(*divisionInfoRequest))
	}).Result((*DivisionShortened)(divisionSample))
	t.SetPanicProcessor(func(err errors.StackTraceError) {
		t.CriticalErr(err)
		t.Quit(err)
	})
	published.Publish(t, t.Quit)
	t.regenerateDatabase()
	return behaviour
}

func (t *TreeEditor) Shutdown() error {
	t.Infoln(t.name, "shut down")
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

func (t *TreeEditor) listFiles() (files, error) {
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

func (t *TreeEditor) getStats() *Statistics {
	stat := &Statistics{
		Positions: len(t.positions),
		Divisions: len(t.divisions)}
	if t.root != nil {
		stat.Companies = len(t.root.Companies)
	}
	return stat
}

func (t *TreeEditor) save(cmd *saveFile) error {
	div, err := t.getCurrent(cmd.id)
	if err != nil {
		return err
	}
	file, err := os.Create(path.Join(t.workDir, path.Base(cmd.file)))
	if err != nil {
		return err
	}
	defer file.Close()
	inspector := &tojson.Inspector{}
	serializer := inspect.NewGenericInspector(inspector)
	div.Inspect(serializer)
	if serializer.GetError() != nil {
		return serializer.GetError()
	}
	_, err = file.Write(inspector.Output())
	return err
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

func (t *TreeEditor) load(cmd *loadFile) (*genericDivision, error) {
	var err error
	t.root, err = LoadFile(path.Join(t.workDir, path.Base(cmd.file)))
	if err != nil {
		return nil, err
	}
	t.regenerateDatabase()
	return (*genericDivision)(t.root.CutTree(0)), nil
}

func (t *TreeEditor) getCurrentNode(cmd *getCurrent) (*genericDivision, error) {
	div := t.findDivisionByPath(cmd.id)
	if div == nil {
		return nil, ErrInvalidId
	}
	return (*genericDivision)(div.CutTree(cmd.depth)), nil
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

func (t *TreeEditor) editPosition(cmd *editPosition) error {
	pos := t.positions[cmd.id]
	if pos.Position == nil {
		return ErrInvalidPosition
	}
	pos.IsEmpty = cmd.isEmpty
	pos.IsSuperior = cmd.isSuperior
	if len(cmd.name) > 0 {
		pos.Name = cmd.name
	}
	if len(cmd.position) > 0 {
		pos.Position.Position = cmd.position
	}
	return nil
}

func (t *TreeEditor) deletePosition(cmd *deletePosition) error {
	pos := t.positions[cmd.id]
	if pos.Position == nil {
		return ErrInvalidPosition
	}
	if pos.parent == nil {
		return ErrPositionIsRoot
	}
	pos.parent.RemovePosition(pos.Position)
	delete(t.positions, pos.Id)
	return nil
}

func (t *TreeEditor) getPosition(cmd *positionInfoRequest) (*Position, error) {
	pos, ok := t.positions[cmd.id]
	if !ok {
		return pos.Position, ErrInvalidPosition
	}
	return pos.Position, nil
}

func (t *TreeEditor) editDivision(cmd *editDivision) error {
	div := t.divisions[cmd.id]
	if div.Division == nil {
		return ErrInvalidDivision
	}
	div.Name = cmd.name
	return nil
}

func (t *TreeEditor) deleteDivision(cmd *deleteDivision) error {
	div := t.divisions[cmd.id]
	if div.Division == nil {
		return ErrInvalidDivision
	}
	if div.parent == nil {
		return ErrCannotRemoveRoot
	}
	t.removeDivisionFromIndex(div.Division)
	div.parent.RemoveDivision(div.Division)
	delete(t.divisions, div.Id)
	return nil
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

func (t *TreeEditor) newDivision(cmd *newDivision) (Id, error) {
	div := t.divisions[cmd.id]
	if div.Division == nil {
		return "", ErrInvalidParent
	}
	newDiv := NewDivision(cmd.name)
	div.AddDivision(newDiv)
	t.divisions[newDiv.Id] = divisionInfo{newDiv, div.Division}
	return Id(newDiv.Id), nil
}

func (t *TreeEditor) newPosition(cmd *newPosition) (Id, error) {
	div := t.divisions[cmd.id]
	if div.Division == nil {
		return "", ErrInvalidParent
	}
	pos := NewPosition(cmd.isSuperior, cmd.name, cmd.position)
	div.AddPosition(pos)
	t.positions[pos.Id] = positionInfo{pos, div.Division}
	return Id(pos.Id), nil
}

func (t *TreeEditor) getDivision(cmd *divisionInfoRequest) (*DivisionShortened, error) {
	div := t.divisions[cmd.id]
	if div.Division == nil {
		return nil, ErrInvalidDivision
	}
	return (*DivisionShortened)(div.Division), nil
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
