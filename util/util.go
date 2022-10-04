package util

import (
	"fmt"
)

func Assert(cond bool, format string, a ...interface{}) {
	if !cond {
		panic(fmt.Errorf(format, a...))
	}
}

func Unreachable() {
	panic("unreachable")
}

func Recover(err *error) {
	if r := recover(); r != nil {
		*err = fmt.Errorf("%v", r)
	}
}
