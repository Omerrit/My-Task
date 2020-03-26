package auth

import (
	"gerrit-share.lan/go/actors"
)

type Filter struct {
	actor  *actors.Actor
	users  IdMap //connection->user map
	handle Handle
}

func (f *Filter) commandFilter(cmd interface{}) error {
	info, ok := cmd.(Info)
	if !ok {
		return nil
	}
	user, ok := f.users[info.ConnectionId()]
	if !ok {
		saved := f.actor.PauseCommand()
		connId := info.ConnectionId()
		f.handle.RequestUserId(info.ConnectionId(), func(userId Id) {
			f.users.Add(connId, userId)
			info.getBase().user = userId
			f.actor.ResumeCommand(saved)
		}, func(err error) {
			f.actor.CancelCommand(saved, err)
		})
		return nil
	}
	info.getBase().user = user
	return nil
}

func (f *Filter) connectionClosed(connId Id) {
	f.users.Delete(connId)
}

func (f *Filter) privilegesChanged(connId Id) {
	f.users.Delete(connId)
}

func (f *Filter) MakeBehaviour(actor *actors.Actor) actors.Behaviour {
	var behaviour actors.Behaviour
	behaviour.PushCommandFilter(f.commandFilter)
	f.actor = actor
	f.handle.Acquire(actor, nil, actor.Quit)
	behaviour.AddMessage(new(connectionClosed), func(msg interface{}) {
		f.connectionClosed(msg.(*connectionClosed).Id)
	}).AddMessage(new(connectionPrivilegesChanged), func(msg interface{}) {
		f.privilegesChanged(msg.(*connectionPrivilegesChanged).Id)
	})
	return behaviour
}
