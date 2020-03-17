package format

import (
	"bytes"
	"gerrit-share.lan/go/tools/sourcecheck/internal/checkers"
	"io/ioutil"
	"os/exec"
)

const gofmt = "gofmt"

type formatChecker struct{}

func NewFormatChecker() checkers.Checker {
	return formatChecker{}
}

func (c formatChecker) Setup(string, string) error { return nil }

func (c formatChecker) Check(fileName string) checkers.Messages {
	formated, err := c.isFormated(fileName)
	var output checkers.Messages
	if err != nil {
		return output.Append(fileName, err.Error(), checkers.NoLine, checkers.NoColumn)
	}
	if !formated {
		return output.Append(fileName, "unformated", checkers.NoLine, checkers.NoColumn)
	}
	return nil
}

func (c formatChecker) callGofmt(filePath string) ([]byte, error) {
	return exec.Command(gofmt, filePath).Output()
}

func (c formatChecker) isFormated(fileName string) (bool, error) {
	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return false, err
	}
	gofmtData, err := c.callGofmt(fileName)
	if err != nil {
		return false, err
	}
	return bytes.Compare(fileData, gofmtData) == 0, nil
}
