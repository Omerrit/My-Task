package basicerrors

import ()

type wrappedError struct {
	error
}

func (w wrappedError) Unwrap() error {
	return w.error
}

type augmentedError struct {
	wrappedError
	basicError BasicError
}

func (a *augmentedError) Is(err error) bool {
	basic, ok := err.(BasicError)
	return ok && basic == a.basicError
}

func (a *augmentedError) As(target interface{}) bool {
	pBasic, ok := target.(*BasicError)
	if !ok {
		return false
	}
	*pBasic = a.basicError
	return true
}

func Augment(err error, code BasicError) error {
	return &augmentedError{wrappedError{err}, code}
}
