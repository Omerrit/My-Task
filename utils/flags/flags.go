package flags

import (
	"flag"
)

func IntFlag(value *int, name string, usage string) {
	flag.IntVar(value, name, *value, usage)
}

func Int64Flag(value *int64, name string, usage string) {
	flag.Int64Var(value, name, *value, usage)
}

func StringFlag(value *string, name string, usage string) {
	flag.StringVar(value, name, *value, usage)
}

func BoolFlag(value *bool, name string, usage string) {
	flag.BoolVar(value, name, *value, usage)
}

func Int32Flag(value *int32, name string, usage string) {
	flag.Var((*int32Value)(value), name, usage)
}
