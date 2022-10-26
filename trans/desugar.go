package trans

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/fun"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
)

func Desugar(expr ast.Expr) ast.Expr {
	switch e := expr.(type) {
	case *ast.StrExpr, *ast.NumExpr, *ast.BoolExpr:
		return expr
	case *ast.TimeExpr:
		return expr
	//	// 'time str' -> strtotime(`time str`)
	//	s := e.Text[1 : len(e.Text)-1]
	//	args := []ast.Expr{ast.Str(fmt.Sprintf("%q", s), s)}
	//	callee := ast.Ident(fun.STRTOTIME)
	//	return ast.Call(callee, args)
	case *ast.ListExpr:
		l := make([]ast.Expr, len(e.Elems))
		for i, el := range e.Elems {
			l[i] = Desugar(el)
		}
		return ast.List(l)
	case *ast.MapExpr:
		m := make([]ast.Pair, len(e.Pairs))
		for i, p := range e.Pairs {
			m[i] = ast.Pair{Key: Desugar(p.Key), Val: Desugar(p.Val)}
		}
		return ast.Map(m)
	case *ast.ObjExpr:
		o := make([]ast.Field, len(e.Fields))
		for i, f := range e.Fields {
			o[i] = ast.Field{Name: f.Name, Val: Desugar(f.Val)}
		}
		return ast.Obj(o)
	case *ast.IdentExpr:
		return expr
	case *ast.UnaryExpr:
		callee := ast.Ident(e.Name)
		args := []ast.Expr{Desugar(e.LHS)}
		return ast.Call(callee, args)
	case *ast.BinaryExpr:
		callee := ast.Ident(e.Name)
		args := []ast.Expr{Desugar(e.LHS), Desugar(e.RHS)}
		return ast.Call(callee, args)
	case *ast.TenaryExpr:
		if e.Name == token.QUESTION {
			// cond ? then : else ~~> if(cond, then, else)
			l := Desugar(e.Left)
			m := Desugar(e.Mid)
			r := Desugar(e.Right)
			// 如果 if 需要处理成特殊语法, 则需要 desugar 成 if-node
			// return ast.If(l, m, r)
			args := []ast.Expr{l, m, r}
			callee := ast.Ident(fun.IF)
			return ast.Call(callee, args)
		}
		util.Unreachable()
		return nil
	case *ast.CallExpr:
		if mem, ok := e.Callee.(*ast.MemberExpr); ok {
			// obj.method(arg...) -> method(obj, arg...)
			args := make([]ast.Expr, len(e.Args)+1)
			args[0] = Desugar(mem.Obj)
			for i, arg := range e.Args {
				args[i+1] = Desugar(arg)
			}
			callee := ast.Ident(mem.Field.Name)
			return ast.Call(callee, args)
		} else {
			args := make([]ast.Expr, len(e.Args))
			for i, arg := range e.Args {
				args[i] = Desugar(arg)
			}
			callee := Desugar(e.Callee)
			return ast.Call(callee, args)
		}
	case *ast.SubscriptExpr:
		return ast.Subscript(Desugar(e.Var), Desugar(e.Idx))
	case *ast.MemberExpr:
		return ast.Member(Desugar(e.Obj), e.Field)
	case *ast.GroupExpr:
		return Desugar(e.SubExpr)
	//case *ast.IfExpr:
	//	// if 是普通函数, 这里不需要
	//	cond := Desugar(e.Cond)
	//	then := Desugar(e.Then)
	//	els := Desugar(e.Else)
	//	// return ast.If(cond, then, els)
	//	callee := ast.Ident(token.IF)
	//	args := []ast.Expr{cond, then, els}
	//	return ast.Call(callee, args)
	default:
		util.Unreachable()
		return nil
	}
}
