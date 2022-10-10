package vm

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/fun"
	"github.com/goghcrow/yae/val"
)

type intrinsicCallByNeed func(*Compiler, *bytecode, []*ast.Expr, *val.Env)

var intrinsicsCallByNeed = map[*val.Val]intrinsicCallByNeed{}
var intrinsicsCallByValue = map[*val.Val]opcode{}

func init() {
	intrinsicsCallByValue[fun.EQ_BOOL_BOOL] = OP_EQ_BOOL_BOOL
	intrinsicsCallByValue[fun.EQ_NUM_NUM] = OP_EQ_NUM_NUM
	intrinsicsCallByValue[fun.EQ_STR_STR] = OP_EQ_STR_STR
	intrinsicsCallByValue[fun.EQ_TIME_TIME] = OP_EQ_TIME_TIME
	intrinsicsCallByValue[fun.EQ_LIST_LIST] = OP_EQ_LIST_LIST
	intrinsicsCallByValue[fun.EQ_MAP_MAP] = OP_EQ_MAP_MAP

	intrinsicsCallByValue[fun.NE_BOOL_BOOL] = OP_NE_BOOL_BOOL
	intrinsicsCallByValue[fun.NE_NUM_NUM] = OP_NE_NUM_NUM
	intrinsicsCallByValue[fun.NE_STR_STR] = OP_NE_STR_STR
	intrinsicsCallByValue[fun.NE_TIME_TIME] = OP_NE_TIME_TIME
	intrinsicsCallByValue[fun.NE_LIST_LIST] = OP_NE_LIST_LIST
	intrinsicsCallByValue[fun.NE_MAP_MAP] = OP_NE_MAP_MAP

	intrinsicsCallByValue[fun.LT_NUM_NUM] = OP_LT_NUM_NUM
	intrinsicsCallByValue[fun.LT_TIME_TIME] = OP_LT_TIME_TIME
	intrinsicsCallByValue[fun.LE_NUM_NUM] = OP_LE_NUM_NUM
	intrinsicsCallByValue[fun.LE_TIME_TIME] = OP_LE_TIME_TIME
	intrinsicsCallByValue[fun.GT_NUM_NUM] = OP_GT_NUM_NUM
	intrinsicsCallByValue[fun.GT_TIME_TIME] = OP_GT_TIME_TIME
	intrinsicsCallByValue[fun.GE_NUM_NUM] = OP_GE_NUM_NUM
	intrinsicsCallByValue[fun.GE_TIME_TIME] = OP_GE_TIME_TIME

	intrinsicsCallByValue[fun.ADD_NUM] = OP_ADD_NUM
	intrinsicsCallByValue[fun.ADD_NUM_NUM] = OP_ADD_NUM_NUM
	intrinsicsCallByValue[fun.ADD_STR_STR] = OP_ADD_STR_STR
	intrinsicsCallByValue[fun.SUB_NUM] = OP_SUB_NUM
	intrinsicsCallByValue[fun.SUB_NUM_NUM] = OP_SUB_NUM_NUM
	intrinsicsCallByValue[fun.SUB_TIME_TIME] = OP_SUB_TIME_TIME
	intrinsicsCallByValue[fun.MUL_NUM_NUM] = OP_MUL_NUM_NUM
	intrinsicsCallByValue[fun.DIV_NUM_NUM] = OP_DIV_NUM_NUM
	intrinsicsCallByValue[fun.MOD_NUM_NUM] = OP_MOD_NUM_NUM
	intrinsicsCallByValue[fun.EXP_NUM_NUM] = OP_EXP_NUM_NUM

	intrinsicsCallByValue[fun.MIN_NUM_NUM] = OP_MIN_NUM_NUM
	intrinsicsCallByValue[fun.MAX_NUM_NUM] = OP_MAX_NUM_NUM

	intrinsicsCallByValue[fun.ABS_NUM] = OP_ABS_NUM
	intrinsicsCallByValue[fun.CEIL_NUM] = OP_CEIL_NUM
	intrinsicsCallByValue[fun.FLOOR_NUM] = OP_FLOOR_NUM
	intrinsicsCallByValue[fun.ROUND_NUM] = OP_ROUND_NUM

	intrinsicsCallByValue[fun.LEN_STR] = OP_LEN_STR
	intrinsicsCallByValue[fun.LEN_LIST] = OP_LEN_LIST
	intrinsicsCallByValue[fun.LEN_MAP] = OP_LEN_MAP

	intrinsicsCallByValue[fun.STRTOTIME_STR] = OP_STRTOTIME_STR

	//intrinsicsCallByValue[fun.STRING_ANY] = OP_STRING_ANY
	//intrinsicsCallByValue[fun.ISSET_MAP_ANY] = OP_ISSET_MAP_ANY
}

func init() {
	intrinsicsCallByNeed[fun.IF_BOOL_ANY_ANY] = func(c *Compiler, b *bytecode, args []*ast.Expr, env *val.Env) {
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
