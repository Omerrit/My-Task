package replies

import "gerrit-share.lan/go/actors"

type StringReply string

func (s StringReply) Visit(visitor actors.ResponseVisitor) {
	visitor.ReplyString(string(s))
}
