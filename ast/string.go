package ast

import (
	"fmt"
	"strings"
)

func (e *Expr) String() string {
	switch e.Type {
	case LITERAL:
		l := e.Literal()
		return l.Val
	case IDENT:
		return e.Ident().Name
	case LIST:
		return fmt.Sprintf("[%s]", stringfyExprs(e.List().Elems))
	case MAP:
		return stringfyMap(e.Map())
	case OBJ:
		return stringfyObj(e.Obj())
	case UNARY:
		u := e.Unary()
		if u.Prefix {
			return fmt.Sprintf("%s%s", u.Type, u.LHS)
		} else {
			return fmt.Sprintf("%s%s", u.LHS, u.Type)
		}
	case BINARY:
		b := e.Binary()
		return fmt.Sprintf("%s %s %s", b.LHS, b.Type.String(), b.RHS)
	case TENARY:
		t := e.Tenary()
		return fmt.Sprintf("%s %s %s %s", t.Left, t.Type.String(), t.Mid, t.Right)
	case CALL:
		c := e.Call()
		return fmt.Sprintf("%s(%s)", c.Callee, stringfyExprs(c.Args))
	case SUBSCRIPT:
		s := e.Subscript()
		return fmt.Sprintf("%s[%s]", s.Var, s.Idx)
	case MEMBER:
		m := e.Member()
		return fmt.Sprintf("%s.%s", m.Obj, m.Field)
	case IF:
		iff := e.If()
		return fmt.Sprintf("if %s then %s else %s", iff.Cond, iff.Then, iff.Else)
	}
	panic("not support exprType")
}

func stringfyMap(m *MapExpr) string {
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

func stringfyObj(m *ObjExpr) string {
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

func stringfyExprs(exprs []*Expr) string {
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
