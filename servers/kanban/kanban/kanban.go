package main

import (
	"gerrit-share.lan/go/actors/starter"
	_ "gerrit-share.lan/go/servers/kanban"
)

func main() {
	starter.Launch("")
}
