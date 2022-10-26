package test

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"testing"
)

func TestRecursive(t *testing.T) {
	lt := types.List(nil).List()
	lt.El = lt.Ty()

	lv := val.List(lt, 0).List()
	lv.Add(lv.Vl())

	t.Log(lv)
}

func TestSortedString(t *testing.T) {
	t.Run("stringify_map", func(t *testing.T) {
		m := val.Map(types.Map(types.Num, types.Num).Map()).Map()
		m.Put(val.Num(2), val.Num(1))
		m.Put(val.Num(1), val.Num(1))
		m.Put(val.Num(3), val.Num(1))
		expect := "[1: 1, 2: 1, 3: 1]"
		actual := m.String()
		if actual != expect {
			t.Errorf("expect %s actual %s", expect, actual)
		}
	})
	t.Run("stringify_obj", func(t *testing.T) {
		o := val.Obj(types.Obj([]types.Field{
			{Name: "b", Val: types.Num},
			{Name: "a", Val: types.Num},
			{Name: "c", Val: types.Num},
		}).Obj()).Obj()
		o.Put("a", val.Num(1))
		o.Put("b", val.Num(1))
		o.Put("c", val.Num(1))
		expect := "{a: 1, b: 1, c: 1}"
		actual := o.String()
		if actual != expect {
			t.Errorf("expect %s actual %s", expect, actual)
		}
	})
}
