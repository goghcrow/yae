package types

import (
	"github.com/goghcrow/yae/util"
)

// IsSubtype k1 <: k2
//func Subtype(k1, k2 *Kind) bool {}

func Equals(x, y *Kind) bool {
	return equals(x, y, util.PtrPtrSet{})
}

func equals(x, y *Kind, inProcess util.PtrPtrSet) bool {
	if x == y {
		return true
	}

	if inProcess.Contains(x, y) {
		return true
	} else {
		inProcess.Add(x, y)
	}

	if x == nil || y == nil {
		return false
	}
	if x.Type != y.Type {
		return false
	}

	if x.Type == TSlot {
		// 这里可以直接比较是因为每次调用 Slot 都会生成唯一值
		return x.Slot().Name == y.Slot().Name
	}

	if x.Type == TMap {
		m1 := x.Map()
		m2 := y.Map()
		return equals(m1.Key, m2.Key, inProcess) && equals(m1.Val, m2.Val, inProcess)
	}

	if x.Type == TTuple {
		return equalsTuple(x.Tuple(), y.Tuple(), inProcess)
	}

	if x.Type == TList {
		return equals(x.List().El, y.List().El, inProcess)
	}

	// structural type (without sequence of fields)
	if x.Type == TObj {
		return equalsObj(x.Obj(), y.Obj(), inProcess)
	}

	if x.Type == TFun {
		return equalsFun(x.Fun(), y.Fun(), inProcess)
	}

	if x.Type == TMaybe {
		return equals(x.Maybe().Elem, y.Maybe().Elem, inProcess)
	}

	return true
}

func equalsObj(x *ObjKind, y *ObjKind, inProcess util.PtrPtrSet) bool {
	if len(x.Fields) != len(y.Fields) {
		return false
	}

	for _, xf := range x.Fields {
		yf, ok := y.GetField(xf.Name)
		if !ok || !equals(xf.Val, yf.Val, inProcess) {
			return false
		}
	}
	return true
}

func equalsTuple(x, y *TupleKind, inProcess util.PtrPtrSet) bool {
	xt := x.Tuple()
	yt := y.Tuple()
	if len(xt.Val) != len(yt.Val) {
		return false
	}
	for i := range xt.Val {
		if !equals(xt.Val[i], yt.Val[i], inProcess) {
			return false
		}
	}
	return true
}

func equalsFun(x, y *FunKind, inProcess util.PtrPtrSet) bool {
	if len(x.Param) != len(y.Param) {
		return false
	}
	for i := range x.Param {
		if !equals(x.Param[i], y.Param[i], inProcess) {
			return false
		}
	}
	return equals(x.Return, y.Return, inProcess)
}
