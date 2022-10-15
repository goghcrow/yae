package util

import "fmt"

func Unreachable() {
	panic("unreachable")
}

func Recover(err *error) {
	if r := recover(); r != nil {
		*err = fmt.Errorf("%v", r)
	}
}
