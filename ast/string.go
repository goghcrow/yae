package ast

import (
	"fmt"
	"github.com/goghcrow/yae/util"
	"strings"
)

func (e *Expr) String() string {
	switch e.Type {
	case LITERAL:
		return e.Literal().Val
	case IDENT:
		return e.Ident().Name
	case LIST:
		return fmt.Sprintf("[%s]", stringifyExprs(e.List().Elems))
	case MAP:
		return stringifyMap(e.Map())
	case OBJ:
		return stringifyObj(e.Obj())
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
		return fmt.Sprintf("%s(%s)", c.Callee, stringifyExprs(c.Args))
	case SUBSCRIPT:
		s := e.Subscript()
		return fmt.Sprintf("%s[%s]", s.Var, s.Idx)
	case MEMBER:
		m := e.Member()
		return fmt.Sprintf("%s.%s", m.Obj, m.Field)
	case IF:
		iff := e.If()
		return fmt.Sprintf("if %s then %s else %s end", iff.Cond, iff.Then, iff.Else)
	default:
		util.Unreachable()
		return ""
	}
}

func stringifyMap(m *MapExpr) string {
	pairs := m.Map().Pairs
	if len(pairs) == 0 {
		return "[:]"
	}
	buf := &strings.Builder{}
	buf.WriteString("[")
	isFst := true
	for _, p := range pairs {
		if isFst {
			isFst = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(p.Key.String())
		buf.WriteString(": ")
		buf.WriteString(p.Val.String())
	}
	buf.WriteString("]")
	return buf.String()
}

func stringifyObj(m *ObjExpr) string {
	fs := m.Obj().Fields
	if len(fs) == 0 {
		return "{}"
	}
	buf := &strings.Builder{}
	buf.WriteString("{")
	isFst := true
	for name, val := range fs {
		if isFst {
			isFst = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(name)
		buf.WriteString(": ")
		buf.WriteString(val.String())
	}
	buf.WriteString("}")
	return buf.String()
}

func stringifyExprs(exprs []*Expr) string {
	buf := &strings.Builder{}
	isFst := true
	for _, elem := range exprs {
		if isFst {
			buf.WriteString(elem.String())
			isFst = false
		} else {
			buf.WriteString(", ")
			buf.WriteString(elem.String())
		}
	}
	return buf.String()
}
