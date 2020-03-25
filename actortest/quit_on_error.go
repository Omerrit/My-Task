package actortest

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/common"
	"testing"
)

func QuitOnError(t *testing.T, actor actors.ActorCompatible, message string) common.ErrorCallback {
	return func(err error) {
		t.Errorf("%s: %s\n", message, err)
		actor.GetBase().Quit(err)
	}
}
