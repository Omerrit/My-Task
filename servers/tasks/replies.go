package tasks

import (
	"database/sql"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"gerrit-share.lan/go/inspect/timetypes"
	"gerrit-share.lan/go/web/auth"
)

type task struct {
	id          auth.Id
	created     timetypes.UnixMilliseconds
	createdBy   auth.Id
	modified    timetypes.UnixMilliseconds
	modifiedBy  auth.Id
	name        string
	description string
	deadline    timetypes.UnixMilliseconds
	status      int32
}

const taskName = packageName + ".task"

func (t *task) Visit(visitor actors.ResponseVisitor) {
	visitor.Reply(t)
}

func (t *task) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(taskName, "detailed task info")
	{
		o.String(&t.name, "name", true, "task name")
		o.String(&t.description, "description", true, "task description")
		t.deadline.Inspect(o.Value("deadline", true, "task deadline"))
		o.Int32(&t.status, "status", true, "task status")
		t.created.Inspect(o.Value("created_on", true, "task creation date"))
		t.createdBy.Inspect(o.Value("created_by", true, "task author"))
		t.modified.Inspect(o.Value("modified_on", true, "last modification date"))
		t.modifiedBy.Inspect(o.Value("modified_by", true, "author of the last modification"))
		o.End()
	}
}

func (t *task) createTable(db *sql.DB) error {
	_, err := db.Exec(`create table tbl_tasks 
	(f_id bit(128) primary key,
	f_name character varying(250),
	f_description text,
	f_deadline timestamp,
	f_status int,
	f_created_on timestamp,
	f_created_by bit(128),
	f_modified_on timestamp,
	f_modified_by bit(128))`)
	return err
}

type taskLoaderById sql.Stmt

func (t *task) newloaderById(db *sql.DB) (*taskLoaderById, error) {
	loader, err := db.Prepare(
		`select 
		f_name,
		f_description,
		f_deadline,
		f_status,
		f_created_on,
		f_created_by,
		f_modified_on,
		f_modified_by from tbl_tasks where f_id=$1`)
	return (*taskLoaderById)(loader), err
}

func (t *taskLoaderById) load(id auth.Id, task *task) (*task, error) {
	err := (*sql.Stmt)(t).QueryRow(id[:]).Scan(
		&task.name,
		&task.description,
		&task.deadline.Time,
		&task.status,
		&task.created.Time,
		&task.createdBy,
		&task.modified.Time,
		&task.modifiedBy)
	task.id = id
	return task, err
}

type taskSaver sql.Stmt

func (t *task) newSaverById(db *sql.DB) (*taskSaver, error) {
	saver, err := db.Prepare(
		`insert into tbl_tasks values 
		(f_id=$1,
		f_name=$2,
		f_description=$3,
		f_status=$4,
		f_created_on=$5,
		f_created_by=$6,
		f_modified_on=$7,
		f_modified_by=$8)`)
	return (*taskSaver)(saver), err
}

func (t *taskSaver) save(task *task) error {
	_, err := (*sql.Stmt)(t).Exec(
		task.id[:],
		task.name,
		task.description,
		task.status,
		task.created.Time,
		task.createdBy[:],
		task.modified.Time,
		task.modifiedBy[:])
	return err
}

func init() {
	inspectables.Register(taskName, func() inspect.Inspectable { return new(task) })
}
