package interp

import (
	"time"

	"github.com/goghcrow/yae/compiler"
	"github.com/goghcrow/yae/parser/ast"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
)

// Interp AST Interpreter
// implement for fun, don't use
func Interp(expr ast.Expr, env *val.Env) compiler.Closure {
	return func(env *val.Env) *val.Val {
		return interp(expr, env)
	}
}

func interp(expr ast.Expr, env *val.Env) *val.Val {
	switch e := expr.(type) {
	case *ast.StrExpr:
		return val.Str(e.Val)

	case *ast.NumExpr:
		return val.Num(e.Val)

	case *ast.TimeExpr:
		return val.Time(time.Unix(e.Val, 0))

	case *ast.BoolExpr:
		if e.Val {
			return val.True
		} else {
			return val.False
		}

	case *ast.ListExpr:
		els := e.Elems
		sz := len(els)
		if sz == 0 {
			// 注意空列表类型 list[nothing]
			ty := types.List(types.Bottom).List()
			return val.List(ty, 0)
		}

		ty := e.Type.(*types.Type).List()
		l := val.List(ty, sz).List()
		for i, el := range els {
			l.V[i] = interp(el, env)
		}
		return l.Vl()

	case *ast.MapExpr:
		sz := len(e.Pairs)
		if sz == 0 {
			// 注意空 map 类型 map[nothing, nothing]
			ty := types.Map(types.Bottom, types.Bottom).Map()
			return val.Map(ty)
		}

		// 保持字面量声明的执行顺序
		ty := e.Type.(*types.Type).Map()
		m := val.Map(ty).Map()
		for _, it := range e.Pairs {
			k := interp(it.Key, env).Key()
			m.V[k] = interp(it.Val, env)
		}
		return m.Vl()

	case *ast.ObjExpr:
		sz := len(e.Fields)
		if sz == 0 {
			ty := types.Obj([]types.Field{}).Obj()
			return val.Obj(ty)
		}

		// 保持字面量声明的执行顺序
		ty := e.Type.(*types.Type).Obj()
		m := val.Obj(ty).Obj()
		for i, f := range e.Fields {
			m.V[i] = interp(f.Val, env)
		}
		return m.Vl()

	case *ast.IdentExpr:
		return env.MustGet(e.Name)

	case *ast.CallExpr:
		fun := resolveFun(e, env)
		args := interpArgs(fun, e, env)
		return fun.Call(args...)

	case *ast.SubscriptExpr:
		// 也可以 desugar 成 build-in-fun
		lhs := interp(e.Var, env)
		rhs := interp(e.Idx, env)

		switch lhs.Type.Kind {
		case types.KList:
			return listSel(lhs, rhs)
		case types.KMap:
			return mapSel(lhs, rhs)
		default:
			util.Unreachable()
			return nil
		}

	case *ast.MemberExpr:
		// 也可以 desugar 成 build-in-fun
		return interp(e.Obj, env).Obj().V[e.Index]

	//case *ast.IfExpr:
	//	// IF 已经 desugar 成 lazyFun 了, 这里已经没用了
	//	// if 分支是 lazy 的 (短路)
	//	if interp(e.Cond, env).Bool().V {
	//		return interp(e.Then, env)
	//	} else {
	//		return interp(e.Else, env)
	//	}

	default:
		util.Unreachable()
		return nil
	}
}

func listSel(lhs, rhs *val.Val) *val.Val {
	lst := lhs.List().V
	idx := int(rhs.Num().V)
	util.Assert(idx < len(lst), "out of range %d of %s", idx, lst)
	return lst[idx]
}

func mapSel(lhs, rhs *val.Val) *val.Val {
	m := lhs.Map()
	k := rhs
	v, ok := m.Get(k)
	// 如果引入 null 、nil 会让类型检查复杂以及做不到 null 安全
	// 可以加一个返回 Maybe 的 get map 函数
	// 或者用 if(isset(m, k), m[k], default)
	util.Assert(ok, "undefined key %s of %s", k, m)
	return v
}

func resolveFun(call *ast.CallExpr, env *val.Env) *val.FunVal {
	if call.Resolved == "" {
		return interp(call.Callee, env).Fun()
	} else {
		if call.Index < 0 {
			return env.MustGetMonoFun(call.Resolved)
		} else {
			return env.MustGetPolyFuns(call.Resolved)[call.Index]
		}
	}
}

func interpArgs(fun *val.FunVal, e *ast.CallExpr, env *val.Env) []*val.Val {
	params := fun.Type.Fun().Param
	args := make([]*val.Val, len(e.Args))
	for i, arg := range e.Args {
		if fun.Lazy {
			// 惰性求值函数参数会被包装成 thunk, 注意没有缓存
			args[i] = thunkify(arg, env, params[i])
		} else {
			args[i] = interp(arg, env)
		}
	}
	return args
}

func thunkify(arg ast.Expr, env *val.Env, retK *types.Type) *val.Val {
	fk := types.Fun("thunk", []*types.Type{}, retK)
	return val.Fun(fk, func(v ...*val.Val) *val.Val {
		return interp(arg, env)
	})
}
