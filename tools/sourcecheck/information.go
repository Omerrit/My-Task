package main

import (
	"fmt"
	"gerrit-share.lan/go/tools/sourcecheck/internal/checkers"
	"path/filepath"
	"strings"
)

type information []info

func (i information) IsEmpty() bool {
	return len(i) == 0
}

func (i information) String() string {
	var buf strings.Builder
	for _, info := range i {
		buf.WriteString(info.String() + "\n")
	}
	return buf.String()
}

type info struct {
	messages          checkers.Messages
	packageImportPath string
}

func (i *info) IsEmpty() bool {
	return i == nil || i.messages.IsEmpty()
}

func (info *info) String() string {
	totalString := ""
	for i, msg := range (*info).messages {
		str := ""
		if info.packageImportPath != "" {
			fileName := filepath.Base(msg.FileName)
			str = fmt.Sprintf("%s/%s", info.packageImportPath, fileName)
		} else {
			str = msg.FileName
		}
		if msg.Line != 0 {
			str = fmt.Sprintf("%s:%v:%v", str, msg.Line, msg.Column)
		}
		str = fmt.Sprintf("%s: %s", str, msg.Description)
		totalString = fmt.Sprintf("%s%s", totalString, str)
		if len(info.messages)-1 != i {
			totalString = fmt.Sprintf("%s\n", totalString)
		}
	}
	return totalString
}
