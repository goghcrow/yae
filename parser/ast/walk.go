package ast

import (
	"github.com/goghcrow/yae/parser/token"
	"github.com/goghcrow/yae/util"
)

type Visitor interface {
	Visit(Expr) Visitor
}

func Walk(v Visitor, expr Expr) {
	if v = v.Visit(expr); v == nil {
		return
	}

	// walk children
	switch e := expr.(type) {
	case *ListExpr:
		for _, el := range e.Elems {
			Walk(v, el)
		}
	case *MapExpr:
		for _, p := range e.Pairs {
			Walk(v, p.Key)
			Walk(v, p.Val)
		}
	case *ObjExpr:
		for _, f := range e.Fields {
			Walk(v, f.Val)
		}
	case *UnaryExpr:
		Walk(v, e.LHS)
	case *BinaryExpr:
		Walk(v, e.LHS)
		Walk(v, e.RHS)
	case *TenaryExpr:
		if e.Name == token.QUESTION {
			Walk(v, e.Left)
			Walk(v, e.Mid)
			Walk(v, e.Right)
		}
		util.Unreachable()
	case *CallExpr:
		if mem, ok := e.Callee.(*MemberExpr); ok {
			Walk(v, mem.Obj)
			for _, arg := range e.Args {
				Walk(v, arg)
			}
		} else {
			for _, arg := range e.Args {
				Walk(v, arg)
			}
			Walk(v, e.Callee)
		}
	case *SubscriptExpr:
		Walk(v, e.Var)
		Walk(v, e.Idx)
	case *MemberExpr:
		Walk(v, e.Obj)
	case *GroupExpr:
		Walk(v, e.SubExpr)
	default:
		util.Unreachable()
	}

	v.Visit(nil)
}
