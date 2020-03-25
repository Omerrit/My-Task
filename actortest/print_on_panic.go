package actortest

import (
	//"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/errors"
	"testing"
)

func PrintOnPanic(t *testing.T, actor *actors.Actor) {
	actor.SetPanicProcessor(func(err errors.StackTraceError) {
		//fmt.Println("panic:", err)
		//fmt.Println(err.StackTrace())
		t.Error("panic: ", err)
		t.Error(err.StackTrace())
		actor.Quit(err)
	})
}
