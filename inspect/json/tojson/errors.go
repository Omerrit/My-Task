package tojson

type toJsonError int

const (
	ErrDuplicateValueCall toJsonError = iota
)

var metadataErrorMessages = map[toJsonError]string{
	ErrDuplicateValueCall: "more than one value inspector call",
}

func (e toJsonError) Error() string {
	return metadataErrorMessages[e]
}
