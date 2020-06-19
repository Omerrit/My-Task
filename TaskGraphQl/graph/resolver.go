package graph

import (
	"gerrit-share.lan/go/graph/model"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	mutex sync.Mutex
	posts []*model.Post
	users map[string] <-chan *model.Post
}
