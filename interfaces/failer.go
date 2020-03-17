package interfaces

import ()

type Failer interface {
	Fail(error)
}
