package compile

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/env"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"strconv"
)

type Closure func(env *env.Env) *val.Val

func (c Closure) String() string { return "Closure" }

// Compile 编译成闭包
// 论文的思路 http://www.iro.umontreal.ca/~feeley/papers/FeeleyLapalmeCL87.pdf
// implement compilers for embedded languages
// 注意区分: env1 是编译期环境, env 是运行时环境
// 能在编译期完成的, 尽可能在编译期计算完
func Compile(env1 *env.Env, expr *ast.Expr) Closure {
	switch expr.Type {
	case ast.LITERAL:
		lit := expr.Literal()
		switch lit.LitType {
		case ast.LIT_STR:
			unquote, err := strconv.Unquote(lit.Val)
			util.Assert(err == nil, "invalid string literal: %s", lit.Val)
			s := val.Str(unquote)
			return func(env *env.Env) *val.Val { return s }
		case ast.LIT_NUM:
			v, _ := util.ParseNum(lit.Val)
			n := val.Num(v)
			return func(env *env.Env) *val.Val { return n }
		case ast.LIT_TRUE:
			return func(env *env.Env) *val.Val { return val.True }
		case ast.LIT_FALSE:
			return func(env *env.Env) *val.Val { return val.False }
		}

	case ast.IDENT:
		id := expr.Ident().Name
		return func(env *env.Env) *val.Val {
			// 如果从性能角度考虑, 所有运行时的符号查找其实都可以从 map 换成 array
			// 1. 需要在编译期把符号 resolve 成数组下标
			// 2. 把运行时环境从 map 展开成数组
			v, _ := env.Get(id)
			return v
		}

	case ast.LIST:
		lst := expr.List()
		els := lst.Elems
		sz := len(els)
		//if sz == 0 {
		//	return func(env *env.Env) *val.Val { return val.List(types.Unit) }
		//}
		// assert sz != 0, typecheck 已经检查过了
		cs := make([]Closure, sz)
		for i, el := range els {
			cs[i] = Compile(env1, el)
		}

		return func(env *env.Env) *val.Val {
			car := cs[0](env)
			l := val.List(car.Kind).List()
			l.V = append(l.V, car)
			for i := 1; i < sz; i++ {
				l.V = append(l.V, cs[i](env))
			}
			return l.Vl()
		}

	case ast.IF:
		iff := expr.If()
		cond := Compile(env1, iff.Cond)
		then := Compile(env1, iff.Then)
		els := Compile(env1, iff.Else)

		// 注意 if 分支是 lazy 的
		return func(env *env.Env) *val.Val {
			if cond(env).Bool().V {
				return then(env)
			} else {
				return els(env)
			}
		}

	case ast.CALL:
		call := expr.Call()
		// 函数在编译期进行链接, 通过 golang 闭包的 upval 传递给运行时
		f, _ := env1.Get(call.Resolved)
		fun := f.Fun()

		sz := len(call.Args)
		cs := make([]Closure, sz)
		for i, arg := range call.Args {
			cs[i] = Compile(env1, arg)
		}

		lazy := fun.Lazy
		retK := fun.Kind.Fun().Return

		return func(env *env.Env) *val.Val {
			args := make([]*val.Val, sz)
			if lazy {
				// 惰性求值函数参数会被包装成 thunk, 注意没有缓存
				for i := 0; i < sz; i++ {
					args[i] = thunkify(cs[i], env, retK)
				}
			} else {
				for i := 0; i < sz; i++ {
					args[i] = cs[i](env)
				}
			}
			return fun.V(args...)
		}

	case ast.SUBSCRIPT:
		sub := expr.Subscript()
		va := Compile(env1, sub.Var)
		idx := Compile(env1, sub.Idx)

		return func(env *env.Env) *val.Val {
			x := va(env)
			if x.Kind.Type == types.TList {
				return x.List().V[int(idx(env).Num().V)]
			} else if x.Kind.Type == types.TMap {
				v, ok := x.Map().Get(idx(env))
				// 这里如果引入 null 、nil 之类又会牵扯到子类型, 实际引入了 bottom 类型
				// 可以业务逻辑里头处理 null, 或者引入一个特殊函数(e.g. isset、has)处理
				// 不能定义成普通函数, 因为要 hack 类型检查
				util.Assert(ok, "undefined key %s of %s", idx(env), x.Map())
				return v
			}
			util.Unreachable()
			return nil
		}

	case ast.MEMBER:
		mem := expr.Member()
		obj := Compile(env1, mem.Obj)
		field := mem.Field.Name

		return func(env *env.Env) *val.Val {
			v, _ := obj(env).Obj().V[field]
			return v
		}

	default:
		util.Unreachable()
	}
	return nil
}

func thunkify(c Closure, env *env.Env, retK *types.Kind) *val.Val {
	fk := types.Fun("thunk", []*types.Kind{}, retK)
	return val.Fun(fk, func(v ...*val.Val) *val.Val {
		return c(env)
	})
}
