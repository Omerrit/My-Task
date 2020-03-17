package replies

import "gerrit-share.lan/go/actors"

type BytesReply []byte

func (b BytesReply) Visit(visitor actors.ResponseVisitor) {
	visitor.ReplyBytes(b)
}

type StringReply string

func (s StringReply) Visit(visitor actors.ResponseVisitor) {
	visitor.ReplyString(string(s))
}
