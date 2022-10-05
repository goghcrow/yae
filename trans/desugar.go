package trans

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
)

func Desugar(expr *ast.Expr) *ast.Expr {
	switch expr.Type {
	case ast.LITERAL:
		lit := expr.Literal()
		if lit.LitType == ast.LIT_TIME {
			// 'time str' ~~> strtotime(`time str`)
			args := []*ast.Expr{ast.LitStr("`" + lit.Text[1:len(lit.Text)-1] + "`")}
			callee := ast.Ident("strtotime")
			return ast.Call(callee, args)
		} else {
			return expr
		}
	case ast.IDENT:
		return expr
	case ast.LIST:
		els := expr.List().Elems
		lst := make([]*ast.Expr, len(els))
		for i, el := range els {
			lst[i] = Desugar(el)
		}
		return ast.List(els)
	case ast.MAP:
		pairs := expr.Map().Pairs
		m := make([]ast.Pair, len(pairs))
		for i, p := range pairs {
			m[i] = ast.Pair{Key: Desugar(p.Key), Val: Desugar(p.Val)}
		}
		return ast.Map(m)
	case ast.OBJ:
		fs := expr.Obj().Fields
		obj := make(map[string]*ast.Expr, len(fs))
		for name, v := range fs {
			obj[name] = Desugar(v)
		}
		return ast.Obj(obj)
	case ast.UNARY:
		u := expr.Unary()
		callee := ast.Ident(u.Name)
		lhs := Desugar(u.LHS)
		args := []*ast.Expr{lhs}
		return ast.Call(callee, args)
	case ast.BINARY:
		b := expr.Binary()
		callee := ast.Ident(b.Name)
		lhs := Desugar(b.LHS)
		rhs := Desugar(b.RHS)
		args := []*ast.Expr{lhs, rhs}
		return ast.Call(callee, args)
	case ast.TENARY:
		t := expr.Tenary()
		if t.Name == token.QUESTION {
			// cond ? then : else ~~> if(cond, then, else)
			l := Desugar(t.Left)
			m := Desugar(t.Mid)
			r := Desugar(t.Right)
			// 如果 if 需要处理成特殊语法, 则需要 desugar 成 if-node
			// return ast.If(l, m, r)
			args := []*ast.Expr{l, m, r}
			callee := ast.Ident(token.IF)
			return ast.Call(callee, args)
		}
		util.Unreachable()
		return nil
	case ast.IF:
		// 如果 if 是普通函数, 这里不需要
		iff := expr.If()
		cond := Desugar(iff.Cond)
		then := Desugar(iff.Then)
		els := Desugar(iff.Else)
		// return ast.If(cond, then, els)
		callee := ast.Ident(token.IF)
		args := []*ast.Expr{cond, then, els}
		return ast.Call(callee, args)
	case ast.CALL:
		call := expr.Call()
		callee := call.Callee
		if callee.Type == ast.MEMBER {
			// obj.method(arg...) ~~> method(obj, arg...)
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
		return nil
	}
}
