package facade

import (
	"fmt"
	"github.com/goghcrow/yae/fun"
	"github.com/goghcrow/yae/trans"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"io"
)

type Expr struct {
	typeCheck *types.Env
	runtime   *val.Env
	trans     []trans.Transform
	dbg       io.Writer
}

func NewExpr() *Expr {
	e := Expr{
		typeCheck: types.NewEnv(),
		runtime:   val.NewEnv(),
		trans:     []trans.Transform{},
	}

	e.initTrans()
	e.initFuns()

	return &e
}

func (e *Expr) EnableDebug(out io.Writer) *Expr {
	e.dbg = out
	return e
}

func (e *Expr) RegisterTransformer(trans trans.Transform) {
	e.trans = append(e.trans, trans)
}

func (e *Expr) RegisterFun(v *val. /*Fun*/ Val) {
	e.typeCheck.RegisterFun(v.Kind)
	e.runtime.RegisterFun(v)
}

func (e *Expr) logf(format string, a ...interface{}) {
	if e.dbg != nil {
		_, _ = fmt.Fprintf(e.dbg, format, a...)
	}
}

func (e *Expr) initTrans() {
	e.RegisterTransformer(trans.Desugar)
}

func (e *Expr) initFuns() {
	e.RegisterFun(fun.IF_BOOL_A_A)

	e.RegisterFun(fun.AND_BOOL_BOOL)
	e.RegisterFun(fun.OR_BOOL_BOOL)
	e.RegisterFun(fun.NOT_BOOL)

	e.RegisterFun(fun.PLUS_NUM)
	e.RegisterFun(fun.PLUS_NUM_NUM)
	e.RegisterFun(fun.PLUS_STR_STR)

	e.RegisterFun(fun.MINUS_NUM)
	e.RegisterFun(fun.MINUS_NUM_NUM)
	e.RegisterFun(fun.MINUS_TIME_TIME)

	e.RegisterFun(fun.MUL_NUM_NUM)

	e.RegisterFun(fun.DIV_NUM_NUM)

	e.RegisterFun(fun.MOD_NUM_NUM)

	e.RegisterFun(fun.EXP_NUM_NUM)

	e.RegisterFun(fun.EQ_BOOL_BOOL)
	e.RegisterFun(fun.EQ_NUM_NUM)
	e.RegisterFun(fun.EQ_STR_STR)
	e.RegisterFun(fun.EQ_TIME_TIME)
	e.RegisterFun(fun.EQ_LIST_LIST)
	e.RegisterFun(fun.EQ_MAP_MAP)

	e.RegisterFun(fun.NE_BOOL_BOOL)
	e.RegisterFun(fun.NE_NUM_NUM)
	e.RegisterFun(fun.NE_STR_STR)
	e.RegisterFun(fun.NE_TIME_TIME)
	e.RegisterFun(fun.NE_LIST_LIST)
	e.RegisterFun(fun.NE_MAP_MAP)

	e.RegisterFun(fun.GT_NUM_NUM)
	e.RegisterFun(fun.GT_TIME_TIME)

	e.RegisterFun(fun.GE_NUM_NUM)
	e.RegisterFun(fun.GE_TIME_TIME)

	e.RegisterFun(fun.LT_NUM_NUM)
	e.RegisterFun(fun.LT_TIME_TIME)

	e.RegisterFun(fun.LE_NUM_NUM)
	e.RegisterFun(fun.LE_TIME_TIME)

	e.RegisterFun(fun.TIME_STR_STR)
	e.RegisterFun(fun.NOW)
	e.RegisterFun(fun.TODAY)
	e.RegisterFun(fun.TODAY_NUM)
	e.RegisterFun(fun.TODAY_STR)

	e.RegisterFun(fun.MAX_NUM_NUM)
	e.RegisterFun(fun.MAX_NUM_NUM_NUM)
	e.RegisterFun(fun.MIN_NUM_NUM)
	e.RegisterFun(fun.MIN_NUM_NUM_NUM)

	e.RegisterFun(fun.LEN_LIST)
	e.RegisterFun(fun.LEN_MAP)
	e.RegisterFun(fun.LEN_STR)

}
