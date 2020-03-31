package actortest

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/errors"
	"testing"
)

func QuitOnError(t *testing.T, actor actors.ActorCompatible, message string) common.ErrorCallback {
	return func(err error) {
		t.Errorf("%s: %s\n", message, errors.FullInfo(err))
		actor.GetBase().Quit(err)
	}
}
