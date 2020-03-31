package errors

type fullInfo struct {
	err error
}

func (f fullInfo) Error() string {
	var ste StackTraceError
	if !As(f.err, &ste) {
		return f.err.Error()
	}
	return ste.Error() + "\n" + ste.StackTrace().String()
}

func (f fullInfo) Unwrap() error {
	return f.err
}

func FullInfo(err error) error {
	return fullInfo{err}
}
