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

// todo 添加新的指令, 特化 fun.BuildIn()

func init() {
	intrinsicsCallByNeed[fun.IF_BOOL_A_A] = func(c *Compiler, b *bytecode, args []*ast.Expr, env *val.Env) {
		b.emitCond(c, args[0], args[1], args[2], env)
	}
	// a && b ~> if (a) b else false
	intrinsicsCallByNeed[fun.LOGIC_AND_BOOL_BOOL] = func(c *Compiler, b *bytecode, args []*ast.Expr, env *val.Env) {
		b.emitCond(c, args[0], args[1], ast.LitFalse(), env)
	}
	// a || b ~> if a true else b
	intrinsicsCallByNeed[fun.LOGIC_OR_BOOL_BOOL] = func(c *Compiler, b *bytecode, args []*ast.Expr, env *val.Env) {
		b.emitCond(c, args[0], ast.LitTrue(), args[1], env)
	}
	intrinsicsCallByNeed[fun.LOGIC_NOT_BOOL] = func(c *Compiler, b *bytecode, args []*ast.Expr, env *val.Env) {
		b.compile(c, args[0], env)
		b.emitOP(OP_LOGICAL_NOT)
	}

	intrinsicsCallByValue[fun.EQ_BOOL_BOOL] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_EQ_BOOL_BOOL) }
	intrinsicsCallByValue[fun.EQ_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_EQ_NUM_NUM) }
	intrinsicsCallByValue[fun.EQ_STR_STR] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_EQ_STR_STR) }
	intrinsicsCallByValue[fun.EQ_TIME_TIME] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_EQ_STR_STR) }
	intrinsicsCallByValue[fun.NE_BOOL_BOOL] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_NE_BOOL_BOOL) }
	intrinsicsCallByValue[fun.NE_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_NE_NUM_NUM) }
	intrinsicsCallByValue[fun.NE_STR_STR] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_NE_STR_STR) }
	intrinsicsCallByValue[fun.NE_TIME_TIME] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_NE_STR_STR) }

	intrinsicsCallByValue[fun.LT_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_LT_NUM_NUM) }
	intrinsicsCallByValue[fun.LT_TIME_TIME] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_LT_TIME_TIME) }
	intrinsicsCallByValue[fun.LE_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_LE_NUM_NUM) }
	intrinsicsCallByValue[fun.LE_TIME_TIME] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_LE_TIME_TIME) }
	intrinsicsCallByValue[fun.GT_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_GT_NUM_NUM) }
	intrinsicsCallByValue[fun.GT_TIME_TIME] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_GT_TIME_TIME) }
	intrinsicsCallByValue[fun.GE_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_GE_NUM_NUM) }
	intrinsicsCallByValue[fun.GE_TIME_TIME] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_GE_TIME_TIME) }

	intrinsicsCallByValue[fun.PLUS_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_ADD_NUM) }
	intrinsicsCallByValue[fun.SUB_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_SUB_NUM) }
	intrinsicsCallByValue[fun.PLUS_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_ADD_NUM_NUM) }
	intrinsicsCallByValue[fun.SUB_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_SUB_NUM_NUM) }
	intrinsicsCallByValue[fun.MUL_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_MUL_NUM_NUM) }
	intrinsicsCallByValue[fun.DIV_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_DIV_NUM_NUM) }
	intrinsicsCallByValue[fun.MOD_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_MOD_NUM_NUM) }
	intrinsicsCallByValue[fun.EXP_NUM_NUM] = func(c *Compiler, b *bytecode, env *val.Env) { b.emitOP(OP_EXP_NUM_NUM) }
}

func (b *bytecode) emitCond(c *Compiler, cond, then, els *ast.Expr, env *val.Env) {
	b.compile(c, cond, env)
	b.emitOP(OP_IF_TRUE)
	emitBranchFalse := b.placeholderForMediumInt()

	b.compile(c, then, env)

	b.emitOP(OP_JUMP)
	emitNext := b.placeholderForMediumInt()

	branchFalse := len(b.code)
	b.compile(c, els, env)

	next := len(b.code)
	emitBranchFalse(branchFalse)
	emitNext(next)
}
