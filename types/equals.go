package types

import (
	"github.com/goghcrow/yae/util"
)

// IsSubtype k1 <: k2
//func Subtype(k1, k2 *Type) bool {}

func Equals(x, y *Type) bool {
	return equals(x, y, util.PtrPtrSet{})
}

func equals(x, y *Type, inProcess util.PtrPtrSet) bool {
	if x == y {
		return true
	}

	if x == nil || y == nil {
		return false
	}

	if x.IsComposite() && y.IsComposite() {
		if inProcess.Contains(x, y) {
			return true
		} else {
			inProcess.Add(x, y)
		}
	}

	if x.Kind != y.Kind {
		return false
	}

	if x.Kind == KTyVar {
		// 这里可以直接比较是因为每次调用 TyVar 都会生成唯一值
		return x.TyVar().Name == y.TyVar().Name
	}

	if x.Kind == KMap {
		m1 := x.Map()
		m2 := y.Map()
		return equals(m1.Key, m2.Key, inProcess) && equals(m1.Val, m2.Val, inProcess)
	}

	if x.Kind == kTuple {
		return equalsTuple(x.Tuple(), y.Tuple(), inProcess)
	}

	if x.Kind == KList {
		return equals(x.List().El, y.List().El, inProcess)
	}

	// structural type (without sequence of fields)
	if x.Kind == KObj {
		return equalsObj(x.Obj(), y.Obj(), inProcess)
	}

	if x.Kind == KFun {
		return equalsFun(x.Fun(), y.Fun(), inProcess)
	}

	if x.Kind == KMaybe {
		return equals(x.Maybe().Elem, y.Maybe().Elem, inProcess)
	}

	return true
}

func equalsObj(x *ObjTy, y *ObjTy, inProcess util.PtrPtrSet) bool {
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

func equalsTuple(x, y *TupleTy, inProcess util.PtrPtrSet) bool {
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

func equalsFun(x, y *FunTy, inProcess util.PtrPtrSet) bool {
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
