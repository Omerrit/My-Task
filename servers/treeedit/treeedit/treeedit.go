package main

import (
	"gerrit-share.lan/go/actors/starter"
	_ "gerrit-share.lan/go/servers/treeedit"
	_ "gerrit-share.lan/go/web/protocols/http/services/httpserver"
)

func main() {
	starter.Launch("")
}
