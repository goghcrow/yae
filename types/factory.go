package types

import (
	"github.com/goghcrow/yae/util"
	"strconv"
)

// Slot 新建 Type Variable
// 📢 每次调用都生成「不相等」的 slot
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
	util.Assert(k.IsPrimitive() || k.Type == TSlot || k.Type == TBottom,
		"invalid type of map's key: %s", k)
	m := MapKind{Kind{TMap}, k, v}
	return &m.Kind
}

func Obj(fields []Field) *Kind {
	t := ObjKind{Kind{TObj}, fields, nil}
	t.Index = make(map[string]int, len(fields))
	for i, f := range fields {
		j, ok := t.Index[f.Name]
		util.Assert(!ok, "duplicated field %s in %d and %d", f.Name, i, j)
		t.Index[f.Name] = i
	}
	return &t.Kind
}

func Fun(name string, param []*Kind, ret *Kind) *Kind {
	t := FunKind{Kind{TFun}, name, param, ret}
	return &t.Kind
}

func Maybe(elem *Kind) *Kind {
	t := MaybeKind{Kind{TMaybe}, elem}
	return &t.Kind
}
