package utils

import (
	"errors"
	"strings"
)

const keySeparator = "\x1e"

type ParsedKey struct {
	Type     string
	Id       string
	PropName string
}

func NewParsedKey(typeName, id, propName string) *ParsedKey {
	return &ParsedKey{typeName, id, propName}
}

func (p *ParsedKey) ToString() string {
	return p.Type + keySeparator + p.Id + keySeparator + p.PropName
}

func ParseKey(key string) (*ParsedKey, error) {
	parts := strings.Split(key, keySeparator)
	length := len(parts)
	if length == 3 {
		return &ParsedKey{
			Type:     parts[length-3],
			Id:       parts[length-2],
			PropName: parts[length-1],
		}, nil
	}
	return nil, errors.New("wrong key format")
}
