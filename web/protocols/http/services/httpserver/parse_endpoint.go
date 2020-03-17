package httpserver

import (
	"bytes"
	//"gerrit-share.lan/go/debug"
	"strings"
)

type ParsedEndpoint struct {
	Name    string
	Method  string
	IsGroup bool
	Depth   int
}

const hierarchySeparator byte = '.'

func getResourceDepth(resource string) int {
	parts := strings.Split(resource, ".")
	return len(parts) - 1
}

func ParseEndpoint(endpoint string) ParsedEndpoint {
	parts := strings.Split(endpoint, string(hierarchySeparator))
	switch len(parts) {
	case 0:
		return ParsedEndpoint{"", "", false, 0}
	case 1:
		return ParsedEndpoint{"", parts[0], false, 0}
	case 2:
		return ParsedEndpoint{parts[1], parts[0], len(parts[1]) == 0, 0}
	}
	if len(parts[len(parts)-1]) != 0 {
		return ParsedEndpoint{strings.Join(parts[1:], string(hierarchySeparator)), parts[0], false, len(parts) - 2}
	}
	return ParsedEndpoint{strings.Join(parts[1:(len(parts)-1)], string(hierarchySeparator)), parts[0], true, len(parts) - 3}
}

type SplitEndpoint struct {
	path    []byte
	indices []int
}

func (s *SplitEndpoint) SetEndpoint(str string, separator byte) {
	var i int
	offset := 0
	s.path = []byte(str)
	if s.path[0] == separator {
		s.path = s.path[1:]
	}
	path := s.path
	for i = bytes.IndexByte(path, separator); i != -1; i = bytes.IndexByte(path, separator) {
		offset += i
		s.indices = append(s.indices, offset)
		s.path[offset] = hierarchySeparator
		path = path[i+1:]
		offset += 1
	}
}

func (s *SplitEndpoint) NumParts() int {
	return len(s.indices) + 1
}

func (s *SplitEndpoint) FullPath() []byte {
	return s.path
}

func (s *SplitEndpoint) CuttedFullPath(front int, back int) []byte {
	return s.path[front:(len(s.path) - back)]
}

func (s *SplitEndpoint) Tail() []byte {
	return s.path[(s.indices[len(s.indices)-1] + 1):]
}

func (s *SplitEndpoint) CurrentPath() []byte {
	return s.path[:(s.indices[len(s.indices)-1])]
}

func (s *SplitEndpoint) Advance() {
	s.indices = s.indices[:(len(s.indices) - 1)]
}

func (s *SplitEndpoint) Cut(targetNumParts int) {
	if targetNumParts > s.NumParts() {
		return
	}
	s.indices = s.indices[:(targetNumParts - 1)]
}

func prepareEndpointName(command string, service string) string {
	parts := strings.Split(command, ".")
	commandName := parts[len(parts)-1]
	return strings.ToLower(commandName + "." + service)
}
