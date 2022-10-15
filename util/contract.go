package util

import (
	"fmt"
)

func Assert(cond bool, format string, a ...interface{}) {
	if !cond {
		panic(fmt.Errorf(format, a...))
	}
}
