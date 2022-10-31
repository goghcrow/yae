package types

import (
	"fmt"
	
	"github.com/goghcrow/yae/util"
)

func (t *Type) String() string {
	return stringify(t, util.PtrSet{})
}

func stringify(ty *Type, inProcess util.PtrSet) string {
	if ty.Kind.IsComposite() {
		if inProcess.Contains(ty) {
			return fmt.Sprintf("recursive-type %s@%p", ty.Kind, ty)
		} else {
			inProcess.Add(ty)
		}
	}

	switch ty.Kind {
	case KNum:
		return "num"
	case KStr:
		return "str"
	case KBool:
		return "bool"
	case KTime:
		return "time"
	case kTuple:
		val := ty.Tuple().Val
		xs := make([]string, len(val))
		for i, el := range val {
			xs[i] = stringify(el, inProcess)
		}
		return util.JoinStr(xs, ", ", "(", ")")
	case KList:
		l := ty.List()
		return fmt.Sprintf("list[%s]", stringify(l.El, inProcess))
	case KMap:
		m := ty.Map()
		return fmt.Sprintf("map[%s, %s]", stringify(m.Key, inProcess), stringify(m.Val, inProcess))
	case KObj:
		fs := ty.Obj().Fields
		xs := make([]string, len(fs))
		for i, f := range fs {
			xs[i] = fmt.Sprintf("%s: %s", f.Name, stringify(f.Val, inProcess))
		}
		return util.JoinStr(xs, ", ", "{", "}")
	case KFun:
		f := ty.Fun()
		xs := make([]string, len(f.Param))
		for i, p := range f.Param {
			xs[i] = stringify(p, inProcess)
		}
		pre := "func " + f.Name + "("
		post := ") " + stringify(f.Return, inProcess)
		return util.JoinStr(xs, ", ", pre, post)
	case KMaybe:
		return fmt.Sprintf("maybe[%s]", stringify(ty.Maybe().Elem, inProcess))
	case KTyVar:
		return "'" + ty.TyVar().Name
	case KTop:
		return "⊤"
	case KBot:
		return "⊥"
	default:
		util.Unreachable()
		return ""
	}
}
