package errors

import "strings"

//import "io"

type ErrorArray []error

func (e ErrorArray) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}
	var builder strings.Builder
	for _, err := range e {
		builder.WriteString("\r\n")
		builder.WriteString(err.Error())
	}
	return builder.String()
}

func (e *ErrorArray) Add(err error) {
	if err != nil {
		switch v := err.(type) {
		case ErrorArray:
			*e = append(*e, v...)
		case *ErrorArray:
			*e = append(*e, (*v)...)
		default:
			*e = append(*e, err)
		}
	}
}

func (e ErrorArray) ToError() error {
	switch len(e) {
	case 0:
		return nil
	case 1:
		return e[0]
	default:
		return e
	}
}

func MakeArray(err ...error) ErrorArray {
	var array ErrorArray
	for _, e := range err {
		array.Add(e)
	}
	return array
}

/*func Close(closable ...io.Closer) error {
	var array ErrorArray
	for _, closer := range closable {
		array.Add(closer.Close())
	}
	return array.ToError()
}*/
