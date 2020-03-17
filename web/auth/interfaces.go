package auth

import ()

type InfoBase struct {
	connId Id
	user   Id
}

func (i *InfoBase) ConnectionId() Id {
	return i.connId
}

func (i *InfoBase) UserId() Id {
	return i.user
}

func (i *InfoBase) getBase() *InfoBase {
	return i
}

type Info interface {
	ConnectionId() Id
	User() Id
	getBase() *InfoBase
}

type WithPath interface {
	Path() string
}
