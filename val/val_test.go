package val

import (
	"github.com/goghcrow/yae/types"
	"testing"
)

func TestRecursive(t *testing.T) {
	lt := types.List(nil).List()
	lt.El = lt.Kd()

	lv := List(lt, 0).List()
	lv.Add(lv.Vl())

	t.Log(lv)
}
