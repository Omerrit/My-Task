package tasks

import (
	"database/sql"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/starter"
	_ "github.com/lib/pq"
	"net/url"
)

type taskService struct {
	actors.Actor
	name string
	db   *sql.DB
}

func (t *taskService) MakeBehaviour() actors.Behaviour {
	//STUB
	return actors.Behaviour{}
}

func newTaskService(name string, postgresUrl url.URL) (*taskService, error) {
	db, err := sql.Open("postgres", postgresUrl.String())
	if err != nil {
		return nil, err
	}
	return &taskService{name: name, db: db}, nil
}

func init() {
	starter.SetCreator(serviceName, func(parent *actors.Actor, name string) (actors.ActorService, error) {
		ts, err := newTaskService(serviceName, postgresUrl)
		if err != nil {
			return nil, err
		}
		ts.DependOn(parent.Service())
		return parent.System().Spawn(ts), nil
	})
}
