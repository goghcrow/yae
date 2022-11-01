package sql

import (
	"fmt"
	"strconv"
	"time"

	"github.com/goghcrow/yae/compiler"
	"github.com/goghcrow/yae/parser/ast"
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
)

const (
	True  = "1"
	False = "0"
)

func Compile(expr ast.Expr, env *val.Env) compiler.Closure {
	return compile(expr, env, 0)
}

func compile(expr ast.Expr, env1 *val.Env, outerPrec oper.BP) compiler.Closure {
	switch e := expr.(type) {
	case *ast.StrExpr:
		x := val.Str(fmtVal(val.Str(e.Val)))
		return func(env *val.Env) *val.Val { return x }
	case *ast.NumExpr:
		x := val.Str(fmtVal(val.Num(e.Val)))
		return func(env *val.Env) *val.Val { return x }
	case *ast.TimeExpr:
		x := val.Str(fmtVal(val.Time(time.Unix(e.Val, 0))))
		return func(env *val.Env) *val.Val { return x }
	case *ast.BoolExpr:
		x := val.Str(fmtVal(val.Bool(e.Val)))
		return func(env *val.Env) *val.Val { return x }
	case *ast.ListExpr:
		sz := len(e.Elems)
		xs := make([]compiler.Closure, sz)
		for i, el := range e.Elems {
			xs[i] = compile(el, env1, outerPrec)
		}
		return func(env *val.Env) *val.Val {
			sxs := make([]string, sz)
			for i := 0; i < sz; i++ {
				sxs[i] = xs[i](env).Str().V
			}
			return val.Str(util.JoinStr(sxs, ", ", "(", ")"))
		}
	case *ast.IdentExpr:
		id := e.Name
		field := val.Str("`" + id + "`")
		return func(env *val.Env) *val.Val {
			v, ok := env.Get(id)
			if ok {
				return val.Str(fmtVal(v))
			} else {
				return field
			}
		}
	case *ast.MemberExpr:
		objId, ok := e.Obj.(*ast.IdentExpr)
		util.Assert(ok, "expect ident actual %s", e.Obj)
		id := objId.Name
		idx := e.Index
		return func(env *val.Env) *val.Val {
			return val.Str(fmtVal(env.MustGet(id).Obj().V[idx]))
		}
	case *ast.CallExpr:
		util.Assert(e.Resolved != "", "only support static dispatch")
		var f *val.FunVal
		if e.Index < 0 {
			f = env1.MustGetMonoFun(e.Resolved)
		} else {
			f = env1.MustGetPolyFuns(e.Resolved)[e.Index]
		}
		// 只有 and or not 需要处理括号
		prec, ok := logicalFunPrecTbl[f.Vl()]
		argc := len(e.Args)
		cs := make([]compiler.Closure, argc)
		for i, arg := range e.Args {
			cs[i] = compile(arg, env1, prec)
		}
		parens := ok && outerPrec > prec
		return func(env *val.Env) *val.Val {
			args := make([]*val.Val, argc)
			for i := 0; i < argc; i++ {
				args[i] = cs[i](env)
			}
			if parens {
				return val.Str("(" + f.Call(args...).Str().V + ")")
			} else {
				return f.Call(args...)
			}
		}
	default:
		util.Assert(false, "unsupported expr: %s", e)
		return nil
	}
}

func fmtVal(v *val.Val) string {
	switch v.Type {
	case types.Bool:
		if v.Bool().V {
			return True
		} else {
			return False
		}
	case types.Num:
		n := v.Num()
		if n.IsInt() {
			return util.FmtInt(n.Int())
		} else {
			return util.FmtFloat(n.V)
		}
	case types.Str:
		return escape(v.Str().V)
	case types.Time:
		return fmt.Sprintf("from_unixtime(%d)", v.Time().V.Unix())
	default:
		util.Assert(false, "unsupported val: %s", v)
		return ""
	}
}

func escape(s string) string {
	return strconv.Quote(s)
}
