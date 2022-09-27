package types

import (
	"github.com/goghcrow/yae/util"
)

// Unify 给出两个类型 A 和 B，找到一组变量替换，
// 使得两者的自由变量经过替换之后可以得到一个相同的类型 C
// 如果 a 和 b 都是 slot 并且 m[a] == m[b]，那么 a b 可以合一，m 不变。
// 如果 a 和 b 都是 primitive 并且相同，那么 a b 可以合一，m 不变。
// 如果 a 是 slot，可以合一，并且需要 m[a] 设置为 b；反之亦然。
// 如果 a 和 b 都是 composite，检查两者的构造器和参数是否都能合一，m 会最多被设置两次。
// 对于其他一切情况，a 和 b 不能合一。
func Unify(x, y *Kind, m map[string]*Kind) *Kind {
	if x.Type == TSlot && y.Type == TSlot && Equals(Subst(x, m), Subst(y, m)) {
		return x
	} else if x.IsPrimitive() && y.IsPrimitive() && x.Type == y.Type {
		return x
	} else if x.IsComposite() && y.IsComposite() && x.Type == y.Type {
		switch x.Type {
		case TList:
			el := Unify(x.List().El, y.List().El, m)
			if el == nil {
				return nil
			}
			return List(el)
		case TMap:
			k := Unify(x.Map().Key, y.Map().Key, m)
			if k == nil {
				return nil
			}
			v := Unify(x.Map().Val, y.Map().Val, m)
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
				u := Unify(xk, yk, m)
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
			fs := make(map[string]*Kind)
			for name, xk := range xfs {
				yk, ok := yfs[name]
				if !ok {
					return nil
				}
				u := Unify(xk, yk, m)
				if u == nil {
					return nil
				}
				fs[name] = u
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
				xp := Subst(xf.Param[i], m)
				yp := Subst(yf.Param[i], m)
				params[i] = Unify(xp, yp, m)
				if params[i] == nil {
					return nil
				}
			}
			ret := Unify(xf.Return, yf.Return, m)
			if ret == nil {
				return nil
			}
			return Fun(xf.Name /**/, params, ret)
		default:
			util.Unreachable()
		}
		return nil
	} else if x.Type == TSlot {
		y1 := Subst(y, m)
		if freeFrom(y1, x.Slot()) {
			k, ok := m[x.Slot().Name]
			if ok && !Equals(k, y1) {
				return nil
			}
			m[x.Slot().Name] = y1 //Subst(y1, m)
			return y1
		} else {
			return nil
		}
	} else if y.Type == TSlot {
		x1 := Subst(x, m)
		if freeFrom(x1, y.Slot()) {
			k, ok := m[y.Slot().Name]
			if ok && !Equals(k, x1) {
				return nil
			}
			m[y.Slot().Name] = x1 //Subst(x1, m)
			return x1
		} else {
			return nil
		}
	} else if y.Type == TBottom {
		return x
	} else if x.Type == TTop {
		return x
	} else {
		return nil
	}
}

// Subst substitution
func Subst(k *Kind, m map[string]*Kind) *Kind {
	switch k.Type {
	case TNum:
		return k
	case TStr:
		return k
	case TBool:
		return k
	case TTime:
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
		return Subst(r, m)
	case TList:
		return List(Subst(k.List().El, m))
	case TMap:
		return Map(
			Subst(k.Map().Key, m),
			Subst(k.Map().Val, m),
		)
	case TTuple:
		t := k.Tuple()
		ks := make([]*Kind, len(t.Val))
		for i, xtv := range t.Val {
			ks[i] = Subst(xtv, m)
		}
		return Tuple(ks)
	case TObj:
		o := k.Obj()
		fs := make(map[string]*Kind, len(o.Fields))
		for name, kind := range o.Fields {
			fs[name] = Subst(kind, m)
		}
		return Obj(fs)
	case TFun:
		f := k.Fun()
		params := make([]*Kind, len(f.Param))
		for i, param := range f.Param {
			params[i] = Subst(param, m)
		}
		return Fun(f.Name, params, Subst(f.Return, m))
	case TTop:
		return k
	case TBottom:
		return k
	default:
		util.Unreachable()
	}
	return nil
}

func freeFrom(k *Kind, s *SlotKind) bool {
	switch k.Type {
	case TNum:
		return true
	case TStr:
		return true
	case TBool:
		return true
	case TTime:
		return true
	case TList:
		return freeFrom(k.List().El, s)
	case TMap:
		return freeFrom(k.Map().Key, s) && freeFrom(k.Map().Val, s)
	case TObj:
		for _, fk := range k.Obj().Fields {
			if !freeFrom(fk, s) {
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
	case TSlot:
		return k.Slot().Name != s.Name
	case TTop:
		return true
	case TBottom:
		return true
	default:
		util.Unreachable()
	}
	return false
}
