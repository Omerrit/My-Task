package endpoints

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"net/http"
)

type handler func(data inspect.Inspectable, writer http.ResponseWriter)

type Endpoint struct {
	handler handler
	creator inspectables.Creator
}

func (e *Endpoint) Handler() handler {
	return e.handler
}

func (e *Endpoint) Creator() inspectables.Creator {
	return e.creator
}

type Endpoints map[string]Endpoint

func (e *Endpoints) Add(path string, handler handler, creator inspectables.Creator) {
	if *e == nil {
		*e = make(Endpoints, 1)
	}
	(*e)[path] = Endpoint{handler, creator}
}
