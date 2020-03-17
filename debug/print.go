package debug

import (
	"fmt"
)

type debugLevel int

const (
	Release debugLevel = iota
	Debug
)

func Println(a ...interface{}) {
	if DebugLevelValue > Release {
		fmt.Println(a...)
	}
}

func Printf(format string, a ...interface{}) {
	if DebugLevelValue > Release {
		fmt.Printf(format, a...)
	}
}
