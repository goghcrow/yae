package types

import (
	"github.com/goghcrow/yae/util"
)

// Unify 给出两个类型 A 和 B, 找到一组变量替换,
// 使得两者的自由变量经过替换之后可以得到一个相同的类型 C
// 如果 a 和 b 都是 slot 并且 m[a] == m[b], 那么 a b 可以合一, m 不变.
// 如果 a 和 b 都是 primitive 并且相同, 那么 a b 可以合一, m 不变.
// 如果 a 是 slot, 可以合一, 并且需要 m[a] 设置为 b；反之亦然.
// 如果 a 和 b 都是 composite, 检查两者的构造器和参数是否都能合一, m 会最多被设置两次.
// 对于其他一切情况, a 和 b 不能合一.
// m 即 substitution, type variable -> type
func Unify(s, t *Type, m map[string]*Type) *Type {
	return unify(s, t, m, util.PtrPtrSet{})
}

func unify(x, y *Type, m map[string]*Type, inProcess util.PtrPtrSet) *Type {
	if x.IsComposite() && y.IsComposite() && inProcess.Contains(x, y) {
		return nil
	} else {
		inProcess.Add(x, y)
	}

	switch {
	case x.Kind == KTyVar && y.Kind == KTyVar && Equals(applySubst(x, m), applySubst(y, m)):
		return x
	case x.IsPrimitive() && y.IsPrimitive() && x.Kind == y.Kind:
		return x
	case x.IsComposite() && y.IsComposite() && x.Kind == y.Kind:
		return unifyComposite(x, y, m, inProcess)
	case x.Kind == KTyVar:
		y1 := applySubst(y, m)
		// x 是 type var, 且 x 没有出现在 y 中
		if freeFrom(y1, x.TyVar()) {
			k, ok := m[x.TyVar().Name]
			if ok && !Equals(k, y1) {
				// 不满足 constrain
				return nil
			}
			m[x.TyVar().Name] = y1 //applySubst(y1, m)
			return y1
		} else {
			return nil
		}
	case y.Kind == KTyVar:
		x1 := applySubst(x, m)
		// y 是 type var, 且 y 没有出现在 x 中
		if freeFrom(x1, y.TyVar()) {
			k, ok := m[y.TyVar().Name]
			if ok && !Equals(k, x1) {
				// 不满足 constrain
				return nil
			}
			m[y.TyVar().Name] = x1 //applySubst(x1, m)
			return x1
		} else {
			return nil
		}
	case y.Kind == KBot:
		return x
	case x.Kind == KTop:
		return x
	default:
		// 不满足 constrain
		return nil
	}
}

func unifyComposite(x, y *Type, m map[string]*Type, inProcess util.PtrPtrSet) *Type {
	switch x.Kind {
	case KList:
		el := unify(x.List().El, y.List().El, m, inProcess)
		if el == nil {
			return nil
		}
		return List(el)
	case KMap:
		k := unify(x.Map().Key, y.Map().Key, m, inProcess)
		if k == nil {
			return nil
		}
		v := unify(x.Map().Val, y.Map().Val, m, inProcess)
		if v == nil {
			return nil
		}
		return Map(k, v)
	case kTuple:
		xtv := x.Tuple().Val
		ytv := y.Tuple().Val
		if len(xtv) != len(ytv) {
			return nil
		}
		ks := make([]*Type, len(xtv))
		for i, xk := range xtv {
			yk := ytv[i]
			u := unify(xk, yk, m, inProcess)
			if u == nil {
				return nil
			}
			ks[i] = u
		}
		return Tuple(ks)
	case KObj:
		xfs := x.Obj().Fields
		yfs := y.Obj().Fields
		if len(xfs) != len(yfs) {
			return nil
		}

		fs := make([]Field, len(xfs))
		for i, xf := range xfs {
			yf, ok := y.Obj().GetField(xf.Name)
			if !ok {
				return nil
			}
			u := unify(xf.Val, yf.Val, m, inProcess)
			if u == nil {
				return nil
			}
			fs[i] = Field{xf.Name, u}
		}
		return Obj(fs)
	case KFun:
		xf := x.Fun()
		yf := y.Fun()
		if len(xf.Param) != len(yf.Param) {
			return nil
		}
		params := make([]*Type, len(xf.Param))
		for i := range xf.Param {
			xp := applySubst(xf.Param[i], m)
			yp := applySubst(yf.Param[i], m)
			params[i] = unify(xp, yp, m, inProcess)
			if params[i] == nil {
				return nil
			}
		}
		ret := unify(xf.Return, yf.Return, m, inProcess)
		if ret == nil {
			return nil
		}
		return Fun(xf.Name /**/, params, ret)
	case KMaybe:
		elm := unify(x.Maybe().Elem, y.Maybe().Elem, m, inProcess)
		if elm == nil {
			return nil
		}
		return Maybe(elm)
	default:
		util.Unreachable()
		return nil
	}
}

// 对 ty 应用 type substitution m
func applySubst(ty *Type, m map[string]*Type) *Type {
	switch ty.Kind {
	case KNum, KStr, KBool, KTime:
		return ty
	case KTyVar:
		r, ok := m[ty.TyVar().Name]
		if !ok {
			return ty
		}
		// 避免递归死循环
		if r.Kind == KTyVar && r.TyVar().Name == ty.TyVar().Name {
			return ty
		}
		return applySubst(r, m)
	case KList:
		return List(applySubst(ty.List().El, m))
	case KMap:
		return Map(
			applySubst(ty.Map().Key, m),
			applySubst(ty.Map().Val, m),
		)
	case kTuple:
		t := ty.Tuple()
		ks := make([]*Type, len(t.Val))
		for i, xtv := range t.Val {
			ks[i] = applySubst(xtv, m)
		}
		return Tuple(ks)
	case KObj:
		o := ty.Obj()
		fs := make([]Field, len(o.Fields))
		for i, f := range o.Fields {
			fs[i] = Field{f.Name, applySubst(f.Val, m)}
		}
		return Obj(fs)
	case KFun:
		f := ty.Fun()
		params := make([]*Type, len(f.Param))
		for i, param := range f.Param {
			params[i] = applySubst(param, m)
		}
		return Fun(f.Name, params, applySubst(f.Return, m))
	case KMaybe:
		return Maybe(applySubst(ty.Maybe().Elem, m))
	case KTop:
		return ty
	case KBot:
		return ty
	default:
		util.Unreachable()
		return nil
	}
}

// occur? s 是否出现在k 中, 这里 free 是指没有出现
func freeFrom(ty *Type, s *TypeVariable) bool {
	switch ty.Kind {
	case KNum, KStr, KBool, KTime, KBot, KTop:
		return true
	case KList:
		return freeFrom(ty.List().El, s)
	case KMap:
		return freeFrom(ty.Map().Key, s) && freeFrom(ty.Map().Val, s)
	case KObj:
		for _, f := range ty.Obj().Fields {
			if !freeFrom(f.Val, s) {
				return false
			}
		}
		return true
	case KFun:
		for _, param := range ty.Fun().Param {
			if !freeFrom(param, s) {
				return false
			}
		}
		return freeFrom(ty.Fun().Return, s)
	case KMaybe:
		return freeFrom(ty.Maybe().Elem, s)
	case KTyVar:
		return ty.TyVar().Name != s.Name
	default:
		util.Unreachable()
		return false
	}
}
