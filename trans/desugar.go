package trans

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
)

func Desugar(expr *ast.Expr) *ast.Expr {
	switch expr.Type {
	case ast.LITERAL:
		return expr
	case ast.IDENT:
		return expr
	case ast.LIST:
		els := expr.List().Elems
		lst := make([]*ast.Expr, len(els))
		for i, el := range els {
			lst[i] = Desugar(el)
		}
		return ast.List(els)
	case ast.UNARY:
		u := expr.Unary()
		callee := ast.Ident(u.Type.Name())
		lhs := Desugar(u.LHS)
		args := []*ast.Expr{lhs}
		return ast.Call(callee, args)
	case ast.BINARY:
		b := expr.Binary()
		callee := ast.Ident(b.Type.Name())
		lhs := Desugar(b.LHS)
		rhs := Desugar(b.RHS)
		args := []*ast.Expr{lhs, rhs}
		return ast.Call(callee, args)
	case ast.TENARY:
		t := expr.Tenary()
		if t.Type == token.QUESTION {
			l := Desugar(t.Left)
			m := Desugar(t.Mid)
			r := Desugar(t.Right)
			// return ast.If(l, m, r)
			args := []*ast.Expr{l, m, r}
			callee := ast.Ident(token.IF.Name())
			return ast.Call(callee, args)
		}
		util.Unreachable()
		return nil
	case ast.IF:
		iff := expr.If()
		cond := Desugar(iff.Cond)
		then := Desugar(iff.Then)
		els := Desugar(iff.Else)
		// return ast.If(cond, then, els)
		args := []*ast.Expr{cond, then, els}
		callee := ast.Ident(token.IF.Name())
		return ast.Call(callee, args)
	case ast.CALL:
		call := expr.Call()
		callee := call.Callee
		if callee.Type == ast.MEMBER {
			mem := callee.Member()
			args := make([]*ast.Expr, len(call.Args)+1)
			args[0] = Desugar(mem.Obj)
			for i, arg := range call.Args {
				args[i+1] = Desugar(arg)
			}
			funName := ast.Ident(mem.Field.Name)
			return ast.Call(funName, args)
		} else {
			callee = Desugar(callee)
			args := make([]*ast.Expr, len(call.Args))
			for i, arg := range call.Args {
				args[i] = Desugar(arg)
			}
			return ast.Call(callee, args)
		}
	case ast.SUBSCRIPT:
		subs := expr.Subscript()
		va := Desugar(subs.Var)
		idx := Desugar(subs.Idx)
		return ast.Subscript(va, idx)
	case ast.MEMBER:
		mem := expr.Member()
		obj := Desugar(mem.Obj)
		return ast.Member(obj, mem.Field)
	default:
		util.Unreachable()
	}

	return nil
}