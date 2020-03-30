package actors

import ()

type actorError int

const (
	ErrActorDead actorError = iota
	ErrActorNull
	ErrCancelled
	ErrUnknownCommand
	ErrNotStreamReply
	ErrBadStream
	ErrStreamConsumerSlow
	ErrWrongTypeRequested
	ErrOffsetOutOfRange
	ErrNotGonnaHappen
	ErrNotFound
	ErrAlreadyRegistered
	ErrApiHandleInvalid
)

var actorErrorMessages = map[actorError]string{
	ErrActorDead:          "actor is dead",
	ErrActorNull:          "trying to send to uninitialized actor handle",
	ErrCancelled:          "request was cancelled",
	ErrUnknownCommand:     "unknown command",
	ErrNotStreamReply:     "'stream request' type reply expected",
	ErrBadStream:          "trying to initialize stream output with uninitialized stream request message",
	ErrStreamConsumerSlow: "stream consumer is too slow",
	ErrWrongTypeRequested: "stream source doesn't support requested type",
	ErrOffsetOutOfRange:   "requested offset is out of range",
	ErrNotGonnaHappen:     "waiting for impossible",
	ErrNotFound:           "not found",
	ErrAlreadyRegistered:  "actor already registered",
	ErrApiHandleInvalid:   "api handle is invalid, try to reinitialize it"}

func (a actorError) Error() string {
	return actorErrorMessages[a]
}
