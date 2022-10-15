package ast

import (
	"fmt"
	"github.com/goghcrow/yae/util"
)

func (e *Expr) String() string {
	switch e.Type {
	case LITERAL:
		return e.Literal().Text
	case IDENT:
		return e.Ident().Name
	case LIST:
		l := e.List()
		xs := make([]string, len(l.Elems))
		for i, el := range l.Elems {
			xs[i] = el.String()
		}
		return util.JoinStr(xs, ", ", "[", "]")
	case MAP:
		m := e.Map()
		pairs := m.Map().Pairs
		if len(pairs) == 0 {
			return "[:]"
		}
		xs := make([]string, len(pairs))
		for i, x := range pairs {
			xs[i] = fmt.Sprintf("%s: %s", x.Key, x.Val)
		}
		return util.JoinStr(xs, ", ", "[", "]")
	case OBJ:
		fs := e.Obj().Fields
		xs := make([]string, len(fs))
		for i, f := range fs {
			xs[i] = fmt.Sprintf("%s: %s", f.Name, f.Val)
		}
		return util.JoinStr(xs, ", ", "{", "}")
	case UNARY:
		u := e.Unary()
		if u.Prefix {
			return fmt.Sprintf("%s%s", u.Name, u.LHS)
		} else {
			return fmt.Sprintf("%s%s", u.LHS, u.Name)
		}
	case BINARY:
		b := e.Binary()
		return fmt.Sprintf("%s %s %s", b.LHS, b.Name, b.RHS)
	case TENARY:
		t := e.Tenary()
		return fmt.Sprintf("%s %s %s %s", t.Left, t.Name, t.Mid, t.Right)
	case CALL:
		c := e.Call()
		xs := make([]string, len(c.Args))
		for i, a := range c.Args {
			xs[i] = a.String()
		}
		return util.JoinStr(xs, ", ", c.Callee.String()+"(", ")")
	case SUBSCRIPT:
		s := e.Subscript()
		return fmt.Sprintf("%s[%s]", s.Var, s.Idx)
	case MEMBER:
		m := e.Member()
		return fmt.Sprintf("%s.%s", m.Obj, m.Field)
	case IF:
		iff := e.If()
		return fmt.Sprintf("if %s then %s else %s end", iff.Cond, iff.Then, iff.Else)
	case GROUP:
		return fmt.Sprintf("(%s)", e.Group().SubExpr)
	default:
		util.Unreachable()
		return ""
	}
}
