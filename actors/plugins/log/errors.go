package log

type logError int

const (
	ErrUnsupportedSeverity logError = iota
)

var logErrorMessages = map[logError]string{
	ErrUnsupportedSeverity: "unsupported log severity"}

func (l logError) Error() string {
	return logErrorMessages[l]
}
