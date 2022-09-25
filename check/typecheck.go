package check

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/env0"
	lex "github.com/goghcrow/yae/lexer"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"strconv"
)

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
			util.Assert(err == nil, "invalid num: %s", lit.Val)
			return types.Num
		case ast.LIT_TRUE:
			v, err := strconv.ParseBool(lit.Val)
			util.Assert(err == nil && v, "invalid bool literal: %s", lit.Val)
			return types.Bool
		case ast.LIT_FALSE:
			v, err := strconv.ParseBool(lit.Val)
			util.Assert(err == nil && !v, "invalid bool literal: %s", lit.Val)
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
		//if sz == 0 {
		//	return types.List(types.Unit)
		//}
		util.Assert(sz != 0, "not support empty list literal")
		elKind := TypeCheck(env0, lst.Elems[0])
		for i := 1; i < sz; i++ {
			kind := TypeCheck(env0, lst.Elems[i])
			util.Assert(types.Equals(elKind, kind),
				"type error, expect %s but %s", elKind, kind)
		}
		return types.List(elKind)

	case ast.IF:
		iff := expr.If()
		condKind := TypeCheck(env0, iff.Cond)
		util.Assert(types.Equals(condKind, types.Bool),
			"type error, expect %s but %s", types.Bool, condKind)
		thenKind := TypeCheck(env0, iff.Then)
		elseKind := TypeCheck(env0, iff.Else)
		util.Assert(types.Equals(thenKind, elseKind),
			"type error, expect %s but %s", thenKind, elseKind)
		return thenKind

	case ast.CALL:
		call := expr.Call()
		util.Assert(call.Callee.Type == ast.IDENT, "invalid fun: %s", call.Callee)
		fn := call.Callee.Ident()

		argSz := len(call.Args)
		args := make([]*types.Kind, argSz)
		for i := 0; i < argSz; i++ {
			args[i] = TypeCheck(env0, call.Args[i])
		}

		if call.Resolved == "" {
			call.Resolved = types.Fun(fn.Name, args, types.Bottom).Fun().OverloadName()
		}
		f, ok := env0.Get(call.Resolved)
		util.Assert(ok, "fun not found: %s", call.Resolved)
		util.Assert(f.Type == types.TFun, "not callable: %s", call.Resolved)
		fun := f.Fun()

		paramSz := len(fun.Param)
		util.Assert(paramSz == argSz, "mismatch arity, expect %d but %d", paramSz, argSz)

		for i := 0; i < paramSz; i++ {
			paramKind := fun.Param[i]
			argKind := TypeCheck(env0, call.Args[i])
			util.Assert(types.Equals(paramKind, argKind), "type error, expect %s but %s", paramKind, argKind)
		}
		return fun.Return

	case ast.SUBSCRIPT:
		sub := expr.Subscript()
		varKind := TypeCheck(env0, sub.Var)

		if varKind.Type == types.TList {
			idxKind := TypeCheck(env0, sub.Idx)
			util.Assert(types.Equals(idxKind, types.Num),
				"type error, expect %s but %s", types.Num, idxKind)
			elKind := varKind.List().El
			return elKind
		}
		if varKind.Type == types.TMap {
			idxKind := TypeCheck(env0, sub.Idx)
			keyKind := varKind.Map().Key
			valKind := varKind.Map().Val
			util.Assert(types.Equals(idxKind, keyKind),
				"type error, expect %s but %s", keyKind, idxKind)
			return valKind
		}
		util.Assert(false, "expect list or map")

	case ast.MEMBER:
		mem := expr.Member()
		objKind := TypeCheck(env0, mem.Obj)
		util.Assert(objKind.Type == types.TObj, "expect %s", types.TObj)
		fn := mem.Field.Name
		fk, ok := objKind.Obj().Fields[fn]
		util.Assert(ok, "undefined filed %s of %s", fn, objKind)
		return fk

	default:
		util.Unreachable()
	}

	return nil
}
