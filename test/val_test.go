package test

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"testing"
)

func TestRecursive(t *testing.T) {
	lt := types.List(nil).List()
	lt.El = lt.Kd()

	lv := val.List(lt, 0).List()
	lv.Add(lv.Vl())

	t.Log(lv)
}
