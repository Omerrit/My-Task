package auth

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

type User struct {
	Id
}

const UserName = packageName + ".user"

func (u *User) Visit(visitor actors.ResponseVisitor) {
	visitor.Reply(u)
}

type userRequest struct {
	Id
}

func NewUser(id Id) *User {
	return &User{id}
}

func init() {
	inspectables.RegisterDescribed(packageName+".user", func() inspect.Inspectable { return new(User) },
		"user id, result of 'request user' command")
	inspectables.RegisterDescribed(packageName+".requser", func() inspect.Inspectable { return new(userRequest) },
		"'request user' command, requires connection id")
}

type getPermission struct {
	connId  Id
	path    string
	command string
}

const getPermissionName = packageName + ".getperm"

func (g *getPermission) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(getPermissionName, "find if this command is permissable for this connection")
	{
		g.connId.Inspect(o.Value("conn_id", true, "connection id"))
		o.String(&g.path, "path", true, "'path' command parameter value")
		o.String(&g.command, "command", true, "command type name")
		o.End()
	}
}

func init() {
	inspectables.Register(getPermissionName, func() inspect.Inspectable { return new(getPermission) })

}
