package ast

import (
	"fmt"
	"github.com/goghcrow/yae/util"
)

func (e *IdentExpr) String() string     { return e.Name }
func (e *StrExpr) String() string       { return e.Text }
func (e *NumExpr) String() string       { return e.Text }
func (e *TimeExpr) String() string      { return e.Text }
func (e *BoolExpr) String() string      { return e.Text }
func (e *SubscriptExpr) String() string { return fmt.Sprintf("%s[%s]", e.Var, e.Idx) }
func (e *MemberExpr) String() string    { return fmt.Sprintf("%s.%s", e.Obj, e.Field) }
func (e *GroupExpr) String() string     { return fmt.Sprintf("(%s)", e.SubExpr) }

func (e *ListExpr) String() string {
	xs := make([]string, len(e.Elems))
	for i, el := range e.Elems {
		xs[i] = el.String()
	}
	return util.JoinStr(xs, ", ", "[", "]")
}

func (e *MapExpr) String() string {
	if len(e.Pairs) == 0 {
		return "[:]"
	}
	xs := make([]string, len(e.Pairs))
	for i, x := range e.Pairs {
		xs[i] = fmt.Sprintf("%s: %s", x.Key, x.Val)
	}
	return util.JoinStr(xs, ", ", "[", "]")
}

func (e *ObjExpr) String() string {
	xs := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		xs[i] = fmt.Sprintf("%s: %s", f.Name, f.Val)
	}
	return util.JoinStr(xs, ", ", "{", "}")
}

func (e *UnaryExpr) String() string {
	if e.Prefix {
		return fmt.Sprintf("%s%s", e.Name, e.LHS)
	} else {
		return fmt.Sprintf("%s%s", e.LHS, e.Name)
	}
}

func (e *BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", e.LHS, e.Name, e.RHS)
}

func (e *TenaryExpr) String() string {
	return fmt.Sprintf("%s %s %s %s", e.Left, e.Name, e.Mid, e.Right)
}

func (e *CallExpr) String() string {
	xs := make([]string, len(e.Args))
	for i, a := range e.Args {
		xs[i] = a.String()
	}
	return util.JoinStr(xs, ", ", e.Callee.String()+"(", ")")
}

//func (e *IfExpr) String() string { return fmt.Sprintf("if %s then %s else %s end", e.Cond, e.Then, e.Else) }
