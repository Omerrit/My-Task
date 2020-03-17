package actors

import ()

type requestCanceller struct {
	id    commandId
	actor *Actor
}

func (r *requestCanceller) Cancel() {
	if r.actor != nil {
		r.actor.cancelRequestById(r.id)
		r.actor = nil
	}
}
