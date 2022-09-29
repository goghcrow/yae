package check

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/env"
	"github.com/goghcrow/yae/env0"
	lex "github.com/goghcrow/yae/lexer"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"strconv"
	"unsafe"
)

// EnvCheck
// env0 compile-env
// env runtime-env
func EnvCheck(env0 *env0.Env, env *env.Env) {
	env0.ForEach(func(name string, kind *types.Kind) {
		v, ok := env.Get(name)
		util.Assert(ok, "undefined %s", name)
		typeAssert(kind, v.Kind, ast.Ident("-env-"))
	})
}

// TypeCheck infer + check
func TypeCheck(env0 *env0.Env, expr *ast.Expr) *types.Kind {
	switch expr.Type {
	case ast.LITERAL:
		lit := expr.Literal()
		switch lit.LitType {
		case ast.LIT_STR:
			return types.Str
		case ast.LIT_NUM:
			_, err := util.ParseNum(lit.Val)
			util.Assert(err == nil, "invalid num literal %s", lit.Val)
			return types.Num
		case ast.LIT_TRUE:
			v, err := strconv.ParseBool(lit.Val)
			util.Assert(err == nil && v, "invalid bool literal %s", lit.Val)
			return types.Bool
		case ast.LIT_FALSE:
			v, err := strconv.ParseBool(lit.Val)
			util.Assert(err == nil && !v, "invalid bool literal %s", lit.Val)
			return types.Bool
		}

	case ast.IDENT:
		id := expr.Ident().Name
		util.Assert(!lex.Reserved(id), "%s reserved", id)
		kind, ok := env0.Get(id)
		util.Assert(ok, "undefined %s", id)
		return kind

	case ast.LIST:
		lst := expr.List()
		sz := len(lst.Elems)
		if sz == 0 {
			return types.List(types.Bottom)
		}

		elKind := TypeCheck(env0, lst.Elems[0])
		for i := 1; i < sz; i++ {
			kind := TypeCheck(env0, lst.Elems[i])
			typeAssert(elKind, kind, expr)
		}
		lst.Kind = types.List(elKind)
		return lst.Kind

	case ast.MAP:
		m := expr.Map()
		sz := len(m.Pairs)
		if sz == 0 {
			return types.Map(types.Bottom, types.Bottom)
		}

		kKind := TypeCheck(env0, m.Pairs[0].Key)
		util.Assert(kKind.IsPrimitive(), "invalid type of map's key: %s", kKind)
		vKind := TypeCheck(env0, m.Pairs[0].Val)
		for i := 1; i < sz; i++ {
			kind := TypeCheck(env0, m.Pairs[i].Key)
			typeAssert(kKind, kind, expr)
			kind = TypeCheck(env0, m.Pairs[i].Val)
			typeAssert(vKind, kind, expr)
		}
		m.Kind = types.Map(kKind, vKind)
		return m.Kind

	case ast.OBJ:
		obj := expr.Obj()
		sz := len(obj.Fields)
		if sz == 0 {
			return types.Obj(map[string]*types.Kind{})
		}

		fs := make(map[string]*types.Kind, sz)
		for name, val := range obj.Fields {
			_, ok := fs[name]
			util.Assert(!ok, "duplicated field %s in %s", name, expr)
			fs[name] = TypeCheck(env0, val)
		}

		obj.Kind = types.Obj(fs)
		return obj.Kind

	/*
		IF 已经 desugar 成 lazyfun 了, 这里不需要
		case ast.IF:
			iff := expr.If()
			condKind := TypeCheck(env0, iff.Cond)
			typeAssert(condKind, types.Bool, expr)
			thenKind := TypeCheck(env0, iff.Then)
			elseKind := TypeCheck(env0, iff.Else)
			typeAssert(thenKind, elseKind, expr)
			return thenKind
	*/

	case ast.CALL:
		call := expr.Call()
		callee := call.Callee

		util.Assert(callee.Type == ast.IDENT, "invalid callable %s in %s", callee, expr)
		fn := callee.Ident()

		argSz := len(call.Args)
		args := make([]*types.Kind, argSz)
		for i := 0; i < argSz; i++ {
			args[i] = TypeCheck(env0, call.Args[i])
		}

		fun := resolve(env0, call, fn.Name, args)

		paramSz := len(fun.Param)
		arityAssert(paramSz, argSz, callee)

		for i := 0; i < paramSz; i++ {
			paramKind := fun.Param[i]
			argKind := TypeCheck(env0, call.Args[i])
			typeAssert(paramKind, argKind, expr)
		}

		call.CalleeKind = fun.Kd()
		return fun.Return

	case ast.SUBSCRIPT:
		sub := expr.Subscript()
		varKind := TypeCheck(env0, sub.Var)

		if varKind.Type == types.TList {
			idxKind := TypeCheck(env0, sub.Idx)
			typeAssert(idxKind, types.Num, expr)
			elKind := varKind.List().El
			return elKind
		}
		if varKind.Type == types.TMap {
			idxKind := TypeCheck(env0, sub.Idx)
			keyKind := varKind.Map().Key
			valKind := varKind.Map().Val
			typeAssert(idxKind, keyKind, expr)
			return valKind
		}
		util.Assert(false, "expect list or map actual %s in %s", varKind, sub.Var)

	case ast.MEMBER:
		mem := expr.Member()
		objKind := TypeCheck(env0, mem.Obj)
		util.Assert(objKind.Type == types.TObj, "expect %s in %s", types.TObj, expr)
		fn := mem.Field.Name
		fk, ok := objKind.Obj().Fields[fn]
		util.Assert(ok, "undefined filed %s of %s in %s", fn, objKind, expr)
		return fk

	default:
		util.Unreachable()
	}

	return nil
}

//goland:noinspection SpellCheckingInspection
func resolve(env0 *env0.Env, call *ast.CallExpr, fnName string, args []*types.Kind) *types.FunKind {
	monofk, mono := types.Fun(fnName, args, types.Bottom /*返回类型无所谓*/).Fun().Lookup()
	util.Assert(mono, "")
	f, ok := env0.Get(monofk)
	if ok {
		util.Assert(f.Type == types.TFun, "non callable of %s in %s", call.Resolved, call)
		call.Resolved = monofk
		monofun := f.Fun()
		return monofun
	}

	polyfk, _ := types.Fun(fnName, args, types.Slot("α")).Fun().Lookup()
	fs, ok := env0.Get(polyfk)
	util.Assert(ok, "undefined fun %s in %s", monofk, call)

	fks := (*[]*types.FunKind)(unsafe.Pointer(fs))
	for i, f := range *fks {
		util.Assert(f.Type == types.TFun, "non callable of %s in %s", fnName, call)
		monof := types.Infer(f, args)
		if monof == nil {
			continue
		}
		call.Resolved = polyfk
		call.Index = i
		return monof
	}
	util.Assert(false, "undefined fun %s in %s", monofk, call)
	return nil
}

func arityAssert(expect, actual int, f *ast.Expr) {
	util.Assert(expect == actual, "%s arity mismatch, expect %d actual %d", f, expect, actual)
}

func typeAssert(expect, actual *types.Kind, f *ast.Expr) {
	eq := types.Equals(expect, actual)
	util.Assert(eq, "type mismatched, expect %s actual %s in %s", expect, actual, f)
}
