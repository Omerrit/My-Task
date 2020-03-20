package auth

import ()

type InfoBase struct {
	connId Id
	user   Id
}

func (i *InfoBase) ConnectionId() Id {
	return i.connId
}

func (i *InfoBase) User() Id {
	return i.user
}

func (i *InfoBase) getBase() *InfoBase {
	return i
}

func (i *InfoBase) SetConnectionId(id Id) {
	i.connId = id
}

type Info interface {
	ConnectionId() Id
	User() Id
	getBase() *InfoBase
	SetConnectionId(id Id)
}

type WithPath interface {
	Path() string
}
