package errors

type stackTrace struct {
	trace CallStack
}

func (s *stackTrace) StackTrace() CallStack {
	return s.trace
}

func (s *stackTrace) replaceStackTrace(newStack CallStack) {
	s.trace = newStack
}

type stringError struct {
	body string
	stackTrace
}

func (e *stringError) Error() string {
	return e.body
}

func (e *stringError) Unwrap() error {
	return nil
}

type wrappedError struct {
	error
	stackTrace
}

func (e *wrappedError) Unwrap() error {
	return e.error
}

type wrappedStringError struct {
	stringError
	err error
}

func (w *wrappedStringError) Unwrap() error {
	return w.err
}

type wrappedStringStackTraceError struct {
	StackTraceError
	body string
}

func (w *wrappedStringStackTraceError) Error() string {
	return w.body
}

func (w *wrappedStringStackTraceError) Unwrap() error {
	return w.StackTraceError
}

type UnknownStackTraceError interface {
	StackTraceError
	Source() interface{}
}

type unknownError struct {
	stringError
	body interface{}
}

func (e *unknownError) Source() interface{} {
	return e.body
}

type described []string

func (d *described) appendDescription(description string) {
	*d = append(*d, description)
}

func (d described) Descriptions() []string {
	return []string(d)
}

type describedError struct {
	StackTraceError
	described
}

func (d *describedError) Error() string {
	var result string
	for _, description := range d.Descriptions() {
		if len(description) > 0 {
			if len(result) > 0 {
				result += ": "
			}
			result += description
		}
	}
	if len(result) == 0 {
		return d.StackTraceError.Error()
	}
	if len(d.StackTraceError.Error()) > 0 {
		result += ": " + d.StackTraceError.Error()
	}
	return result
}
