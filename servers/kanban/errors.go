package kanban

import ()

type serverError int

const (
	ErrIdIsNotEmpty serverError = iota
	ErrIdIsInUse
)

var serverErrorStrings = map[serverError]string{
	ErrIdIsNotEmpty: "id have data associated with it, can't delete",
	ErrIdIsInUse:    "id is already in use"}

func (s serverError) Error() string {
	return serverErrorStrings[s]
}
