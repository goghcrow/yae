package trans

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/fun"
	"github.com/goghcrow/yae/loc"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
)

func Desugar(expr ast.Expr) ast.Expr {
	switch e := expr.(type) {
	case *ast.StrExpr, *ast.NumExpr, *ast.BoolExpr:
		return expr
	case *ast.TimeExpr:
		return expr
		//// 'time str' -> strtotime(`time str`)
		//s := e.Text[1 : len(e.Text)-1]
		//args := []ast.Expr{ast.Str(strconv.Quote(s))}
		//callee := ast.Var(fun.STRTOTIME)
		//return ast.Call(callee, args)
	case *ast.ListExpr:
		l := make([]ast.Expr, len(e.Elems))
		for i, el := range e.Elems {
			l[i] = Desugar(el)
		}
		return ast.List(l, e.Loc)
	case *ast.MapExpr:
		m := make([]ast.Pair, len(e.Pairs))
		for i, p := range e.Pairs {
			m[i] = ast.Pair{Key: Desugar(p.Key), Val: Desugar(p.Val)}
		}
		return ast.Map(m, e.Loc)
	case *ast.ObjExpr:
		o := make([]ast.Field, len(e.Fields))
		for i, f := range e.Fields {
			o[i] = ast.Field{Name: f.Name, Val: Desugar(f.Val)}
		}
		return ast.Obj(o, e.Loc)
	case *ast.IdentExpr:
		return expr
	case *ast.UnaryExpr:
		callee := ast.Var(e.Name, e.IdentExpr.Loc)
		args := []ast.Expr{Desugar(e.LHS)}
		dbgCol := loc.DbgCol(e.IdentExpr.Col)
		return ast.Call(callee, args, dbgCol, e.Loc)
	case *ast.BinaryExpr:
		callee := ast.Var(e.Name, e.IdentExpr.Loc)
		args := []ast.Expr{Desugar(e.LHS), Desugar(e.RHS)}
		dbgCol := loc.DbgCol(e.IdentExpr.Col)
		return ast.Call(callee, args, dbgCol, e.Loc)
	case *ast.TenaryExpr:
		if e.Name == token.QUESTION {
			// cond ? then : else ~~> if(cond, then, else)
			l := Desugar(e.Left)
			m := Desugar(e.Mid)
			r := Desugar(e.Right)
			// 如果 if 需要处理成特殊语法, 则需要 desugar 成 if-node
			// return ast.If(l, m, r)
			args := []ast.Expr{l, m, r}
			callee := ast.Var(fun.IF, e.IdentExpr.Loc)
			dbgCol := loc.DbgCol(e.IdentExpr.Col)
			return ast.Call(callee, args, dbgCol, e.Loc)
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
			callee := ast.Var(mem.Field.Name, mem.Field.Loc)
			return ast.Call(callee, args, e.DbgCol, e.Loc)
		} else {
			args := make([]ast.Expr, len(e.Args))
			for i, arg := range e.Args {
				args[i] = Desugar(arg)
			}
			callee := Desugar(e.Callee)
			return ast.Call(callee, args, e.DbgCol, e.Loc)
		}
	case *ast.SubscriptExpr:
		return ast.Subscript(Desugar(e.Var), Desugar(e.Idx), e.DbgCol, e.Loc)
	case *ast.MemberExpr:
		return ast.Member(Desugar(e.Obj), e.Field, e.DbgCol, e.Loc)
	case *ast.GroupExpr:
		return Desugar(e.SubExpr)
	//case *ast.IfExpr:
	//	// if 是普通函数, 这里不需要
	//	cond := Desugar(e.Cond)
	//	then := Desugar(e.Then)
	//	els := Desugar(e.Else)
	//	// return ast.If(cond, then, els)
	//	callee := ast.Var(token.IF)
	//	args := []ast.Expr{cond, then, els}
	//	return ast.Call(callee, args)
	default:
		util.Unreachable()
		return nil
	}
}
