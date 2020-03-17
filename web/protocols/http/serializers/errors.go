package serializers

import "gerrit-share.lan/go/basicerrors"

type urlError int

const (
	ErrUnsupportedType urlError = iota
)

var metadataErrorMessages = map[urlError]string{
	ErrUnsupportedType: "unsupported type"}

func (p urlError) Error() string {
	return metadataErrorMessages[p]
}

func (p urlError) Unwrap() error {
	return basicerrors.BadParameter
}
