package facade

import (
	"fmt"
	"github.com/goghcrow/yae/env"
	"github.com/goghcrow/yae/env0"
	"github.com/goghcrow/yae/fun"
	"github.com/goghcrow/yae/trans"
	"github.com/goghcrow/yae/val"
	"io"
)

type Expr struct {
	typeCheck *env0.Env
	runtime   *env.Env
	trans     []trans.Transform
	dbg       io.Writer
}

func NewExpr() *Expr {
	e := Expr{
		typeCheck: env0.NewEnv(),
		runtime:   env.NewEnv(),
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

func (e *Expr) RegisterFun(v *val.FunVal) {
	e.typeCheck.RegisterFun(v.Kind.Fun())
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
	e.RegisterFun(fun.AND_BOOL_BOOL.Fun())
	e.RegisterFun(fun.OR_BOOL_BOOL.Fun())
	e.RegisterFun(fun.NOT_BOOL.Fun())

	e.RegisterFun(fun.PLUS_NUM.Fun())
	e.RegisterFun(fun.PLUS_NUM_NUM.Fun())
	e.RegisterFun(fun.PLUS_STR_STR.Fun())

	e.RegisterFun(fun.MINUS_NUM.Fun())
	e.RegisterFun(fun.MINUS_NUM_NUM.Fun())
	e.RegisterFun(fun.MINUS_TIME_TIME.Fun())

	e.RegisterFun(fun.MUL_NUM_NUM.Fun())

	e.RegisterFun(fun.DIV_NUM_NUM.Fun())

	e.RegisterFun(fun.MOD_NUM_NUM.Fun())

	e.RegisterFun(fun.EXP_NUM_NUM.Fun())

	e.RegisterFun(fun.EQ_BOOL_BOOL.Fun())
	e.RegisterFun(fun.EQ_NUM_NUM.Fun())
	e.RegisterFun(fun.EQ_STR_STR.Fun())
	e.RegisterFun(fun.EQ_TIME_TIME.Fun())
	e.RegisterFun(fun.EQ_LIST_LIST.Fun())
	e.RegisterFun(fun.EQ_MAP_MAP.Fun())
	e.RegisterFun(fun.EQ_OBJ_OBJ.Fun())

	e.RegisterFun(fun.NE_BOOL_BOOL.Fun())
	e.RegisterFun(fun.NE_NUM_NUM.Fun())
	e.RegisterFun(fun.NE_STR_STR.Fun())
	e.RegisterFun(fun.NE_TIME_TIME.Fun())
	e.RegisterFun(fun.NE_LIST_LIST.Fun())
	e.RegisterFun(fun.NE_MAP_MAP.Fun())
	e.RegisterFun(fun.NE_OBJ_OBJ.Fun())

	e.RegisterFun(fun.GT_NUM_NUM.Fun())
	e.RegisterFun(fun.GT_TIME_TIME.Fun())

	e.RegisterFun(fun.GE_NUM_NUM.Fun())
	e.RegisterFun(fun.GE_TIME_TIME.Fun())

	e.RegisterFun(fun.LT_NUM_NUM.Fun())
	e.RegisterFun(fun.LT_TIME_TIME.Fun())

	e.RegisterFun(fun.LE_NUM_NUM.Fun())
	e.RegisterFun(fun.LE_TIME_TIME.Fun())

	e.RegisterFun(fun.TIME_STR_STR.Fun())
	e.RegisterFun(fun.NOW.Fun())
	e.RegisterFun(fun.TODAY.Fun())
	e.RegisterFun(fun.TODAY_NUM.Fun())

	e.RegisterFun(fun.MAX_NUM_NUM.Fun())
	e.RegisterFun(fun.MAX_NUM_NUM_NUM.Fun())
	e.RegisterFun(fun.MIN_NUM_NUM.Fun())
	e.RegisterFun(fun.MIN_NUM_NUM_NUM.Fun())

}
