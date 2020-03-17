package checkers

const (
	NoLine   = 0
	NoColumn = 0
)

type message struct {
	Column      int
	Line        int
	FileName    string
	Description string
}

type Messages []message

func (m *Messages) Append(fileName string, description string, line int, column int) Messages {
	if *m == nil {
		*m = make(Messages, 0)
	}
	*m = append(*m, message{
		Line:        line,
		Column:      column,
		FileName:    fileName,
		Description: description,
	})
	return *m
}

func (m *Messages) AppendMessages(other Messages) Messages {
	*m = append(*m, other...)
	return *m
}

func (m *Messages) IsEmpty() bool {
	return m == nil || len(*m) == 0
}
