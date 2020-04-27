package utils

import (
	"strings"
)

const keySeparator = "\x1e"

type ParsedKey struct {
	Type     string
	Id       string
	PropName string
}

func IsKeyValid(key string) bool {
	return len(strings.Split(key, keySeparator)) >= 3
}

func ParseKey(key string) *ParsedKey {
	parts := strings.Split(key, keySeparator)
	length := len(parts)
	return &ParsedKey{
		Type:     parts[length-3],
		Id:       parts[length-2],
		PropName: parts[length-1],
	}
}
