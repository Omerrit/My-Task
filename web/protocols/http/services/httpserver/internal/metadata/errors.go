package metadata

type metadataError int

const (
	ErrDuplicateCallValue metadataError = iota
	ErrSameTypeName
)

var metadataErrorMessages = map[metadataError]string{
	ErrDuplicateCallValue: "duplicate call in value inspector",
	ErrSameTypeName:       "types with same name, but different description have been detected"}

func (p metadataError) Error() string {
	return metadataErrorMessages[p]
}
