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
func Unify(s, t *Kind, m map[string]*Kind) *Kind {
	return unify(s, t, m, util.PtrPtrSet{})
}

func unify(x, y *Kind, m map[string]*Kind, inProcess util.PtrPtrSet) *Kind {
	if x.IsComposite() && y.IsComposite() && inProcess.Contains(x, y) {
		return nil
	} else {
		inProcess.Add(x, y)
	}

	switch {
	case x.Type == TSlot && y.Type == TSlot && Equals(applySubst(x, m), applySubst(y, m)):
		return x
	case x.IsPrimitive() && y.IsPrimitive() && x.Type == y.Type:
		return x
	case x.IsComposite() && y.IsComposite() && x.Type == y.Type:
		return unifyComposite(x, y, m, inProcess)
	case x.Type == TSlot:
		y1 := applySubst(y, m)
		// x 是 type var, 且 x 没有出现在 y 中
		if freeFrom(y1, x.Slot()) {
			k, ok := m[x.Slot().Name]
			if ok && !Equals(k, y1) {
				// 不满足 constrain
				return nil
			}
			m[x.Slot().Name] = y1 //applySubst(y1, m)
			return y1
		} else {
			return nil
		}
	case y.Type == TSlot:
		x1 := applySubst(x, m)
		// y 是 type var, 且 y 没有出现在 x 中
		if freeFrom(x1, y.Slot()) {
			k, ok := m[y.Slot().Name]
			if ok && !Equals(k, x1) {
				// 不满足 constrain
				return nil
			}
			m[y.Slot().Name] = x1 //applySubst(x1, m)
			return x1
		} else {
			return nil
		}
	case y.Type == TBottom:
		return x
	case x.Type == TTop:
		return x
	default:
		// 不满足 constrain
		return nil
	}
}

func unifyComposite(x, y *Kind, m map[string]*Kind, inProcess util.PtrPtrSet) *Kind {
	switch x.Type {
	case TList:
		el := unify(x.List().El, y.List().El, m, inProcess)
		if el == nil {
			return nil
		}
		return List(el)
	case TMap:
		k := unify(x.Map().Key, y.Map().Key, m, inProcess)
		if k == nil {
			return nil
		}
		v := unify(x.Map().Val, y.Map().Val, m, inProcess)
		if v == nil {
			return nil
		}
		return Map(k, v)
	case TTuple:
		xtv := x.Tuple().Val
		ytv := y.Tuple().Val
		if len(xtv) != len(ytv) {
			return nil
		}
		ks := make([]*Kind, len(xtv))
		for i, xk := range xtv {
			yk := ytv[i]
			u := unify(xk, yk, m, inProcess)
			if u == nil {
				return nil
			}
			ks[i] = u
		}
		return Tuple(ks)
	case TObj:
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
	case TFun:
		xf := x.Fun()
		yf := y.Fun()
		if len(xf.Param) != len(yf.Param) {
			return nil
		}
		params := make([]*Kind, len(xf.Param))
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
	case TMaybe:
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

// 对 term k 应用 type substitution m
func applySubst(k *Kind, m map[string]*Kind) *Kind {
	switch k.Type {
	case TNum, TStr, TBool, TTime:
		return k
	case TSlot:
		r, ok := m[k.Slot().Name]
		if !ok {
			return k
		}
		// 避免递归死循环
		if r.Type == TSlot && r.Slot().Name == k.Slot().Name {
			return k
		}
		return applySubst(r, m)
	case TList:
		return List(applySubst(k.List().El, m))
	case TMap:
		return Map(
			applySubst(k.Map().Key, m),
			applySubst(k.Map().Val, m),
		)
	case TTuple:
		t := k.Tuple()
		ks := make([]*Kind, len(t.Val))
		for i, xtv := range t.Val {
			ks[i] = applySubst(xtv, m)
		}
		return Tuple(ks)
	case TObj:
		o := k.Obj()
		fs := make([]Field, len(o.Fields))
		for i, f := range o.Fields {
			fs[i] = Field{f.Name, applySubst(f.Val, m)}
		}
		return Obj(fs)
	case TFun:
		f := k.Fun()
		params := make([]*Kind, len(f.Param))
		for i, param := range f.Param {
			params[i] = applySubst(param, m)
		}
		return Fun(f.Name, params, applySubst(f.Return, m))
	case TMaybe:
		return Maybe(applySubst(k.Maybe().Elem, m))
	case TTop:
		return k
	case TBottom:
		return k
	default:
		util.Unreachable()
		return nil
	}
}

// occur? s 是否出现在k 中, 这里 free 是指没有出现
func freeFrom(k *Kind, s *SlotKind) bool {
	switch k.Type {
	case TNum, TStr, TBool, TTime, TBottom, TTop:
		return true
	case TList:
		return freeFrom(k.List().El, s)
	case TMap:
		return freeFrom(k.Map().Key, s) && freeFrom(k.Map().Val, s)
	case TObj:
		for _, f := range k.Obj().Fields {
			if !freeFrom(f.Val, s) {
				return false
			}
		}
		return true
	case TFun:
		for _, param := range k.Fun().Param {
			if !freeFrom(param, s) {
				return false
			}
		}
		return freeFrom(k.Fun().Return, s)
	case TMaybe:
		return freeFrom(k.Maybe().Elem, s)
	case TSlot:
		return k.Slot().Name != s.Name
	default:
		util.Unreachable()
		return false
	}
}
