package types

import (
	"github.com/goghcrow/yae/util"
	"strconv"
)

// Slot 每次调用都生成不相等的 slot
var Slot = func() func(name string) *Kind {
	n := 0
	return func(name string) *Kind {
		n++
		t := SlotKind{Kind{TSlot}, name + strconv.Itoa(n)}
		return &t.Kind
	}
}()

func Tuple(val []*Kind) *Kind {
	t := TupleKind{Kind{TTuple}, val}
	return &t.Kind
}

func List(el *Kind) *Kind {
	t := ListKind{Kind{TList}, el}
	return &t.Kind
}

func Map(k, v *Kind) *Kind {
	util.Assert(k.IsPrimitive() || k.Type == TSlot, "invalid type of map's key: %s", k)
	m := MapKind{Kind{TMap}, k, v}
	return &m.Kind
}

func Obj(fields map[string]*Kind) *Kind {
	t := ObjKind{Kind{TObj}, fields}
	return &t.Kind
}

func Fun(name string, param []*Kind, ret *Kind) *Kind {
	t := FunKind{Kind{TFun}, name, param, ret}
	return &t.Kind
}
