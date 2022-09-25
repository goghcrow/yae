package ast

import "fmt"

func (e *Expr) String() string {
	switch e.Type {
	case LITERAL:
		l := e.Literal()
		return l.Val
	case IDENT:
		return e.Ident().Name
	case LIST:
		return fmt.Sprintf("[%s]", stringExprs(e.List().Elems))
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
		return fmt.Sprintf("%s(%s)", c.Callee, stringExprs(c.Args))
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

func stringExprs(exprs []*Expr) string {
	buf := ""
	isFst := true
	for _, elem := range exprs {
		if isFst {
			buf += elem.String()
			isFst = false
		} else {
			buf += ", " + elem.String()
		}
	}
	return buf
}
