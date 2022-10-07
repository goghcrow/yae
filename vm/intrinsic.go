package vm

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/fun"
	"github.com/goghcrow/yae/val"
)

type intrinsicCallByNeed func(*Compiler, *bytecode, []*ast.Expr, *val.Env)
type intrinsicCallByValue func(*Compiler, *bytecode, *val.Env)

var intrinsicsCallByNeed = map[*val.Val]intrinsicCallByNeed{}
var intrinsicsCallByValue = map[*val.Val]intrinsicCallByValue{}

// todo 添加新的指令, 特化 fun.BuildIn() 中大部分数学和逻辑函数,

func init() {
	intrinsicsCallByNeed[fun.IF_BOOL_A_A] = func(c *Compiler, b *bytecode, args []*ast.Expr, env *val.Env) {
		// todo 重构下这里, 太丑陋了...
		cond := args[0]
		then := args[1]
		els := args[2]

		b.compile(c, cond, env)

		b.emitOP(OP_IF)

		fixBranchTrue := len(b.code)
		branchTrue := 0
		b.emitMediumInt(branchTrue)

		fixBranchFalse := len(b.code)
		branchFalse := 0
		b.emitMediumInt(branchFalse)

		branchTrue = len(b.code)
		b.compile(c, then, env)
		b.emitOP(OP_JUMP)

		fixNxt := len(b.code)
		nxt := 0
		b.emitMediumInt(nxt)

		branchFalse = len(b.code)
		b.compile(c, els, env)

		nxt = len(b.code)

		b.setMediumInt(fixBranchTrue, branchTrue)
		b.setMediumInt(fixBranchFalse, branchFalse)
		b.setMediumInt(fixNxt, nxt)
	}
	//intrinsicsCallByNeed[fun.LOGIC_AND_BOOL_BOOL] = func(c *Compiler, b *bytecode, args []*ast.Expr, env *val.Env) {
	//
	//}
	//intrinsicsCallByNeed[fun.LOGIC_OR_BOOL_BOOL] = func(c *Compiler, b *bytecode, args []*ast.Expr, env *val.Env) {
	//
	//}
	//intrinsicsCallByNeed[fun.LOGIC_NOT_BOOL] = func(c *Compiler, b *bytecode, args []*ast.Expr, env *val.Env) {
	//
	//}

	intrinsicsCallByValue[fun.PLUS_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_ADD_NUM) }
	intrinsicsCallByValue[fun.PLUS_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_ADD_NUM_NUM) }
	intrinsicsCallByValue[fun.SUB_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_SUB_NUM) }
	intrinsicsCallByValue[fun.SUB_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_SUB_NUM_NUM) }
	intrinsicsCallByValue[fun.MUL_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_MUL_NUM_NUM) }
	intrinsicsCallByValue[fun.DIV_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_DIV_NUM_NUM) }
	intrinsicsCallByValue[fun.MOD_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_MOD_NUM_NUM) }
	intrinsicsCallByValue[fun.EXP_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_EXP_NUM_NUM) }
}
