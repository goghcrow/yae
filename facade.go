package expr

import (
	"errors"
	"fmt"
	"github.com/goghcrow/yae/compile"
	"github.com/goghcrow/yae/conv"
	"github.com/goghcrow/yae/fun"
	lex "github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/parser"
	"github.com/goghcrow/yae/trans"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"io"
	"runtime/debug"
)

// Eval without cache compiled closure
func Eval(input string, v interface{}) (*val.Val, error) {
	expr := NewExpr() //.EnableDebug(os.Stderr)
	compileTimeEnv, err := conv.TypeEnvOf(v)
	if err != nil {
		return nil, err
	}
	closure, err := expr.Compile(input, compileTimeEnv)
	if err != nil {
		return nil, err
	}

	runtimeEnv, err := conv.ValEnvOf(v)
	if err != nil {
		return nil, err
	}
	return closure(runtimeEnv)
}

type Expr struct {
	typeCheck *types.Env
	runtime   *val.Env
	trans     []trans.Transform
	dbg       io.Writer
}

type Compiled func(env1 *val.Env) (*val.Val, error)

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

func (e *Expr) Compile(expr string, env0 *types.Env) (c Compiled, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e.dbg != nil {
				e.logf("COMPILE FAIL:\n")
				_, _ = e.dbg.Write(debug.Stack())
			}
			err = errors.New(fmt.Sprintf("%s", r))
		}
	}()

	closure := e.steps(expr, env0)
	c = e.makeCompiled(closure, env0)
	return
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
	for _, f := range fun.Funs {
		e.RegisterFun(f)
	}
}

func (e *Expr) steps(expr string, env0 *types.Env) compile.Closure {
	e.logf("expr: %s\n", expr)

	toks := lex.Lex(expr)
	e.logf("lexed: %s\n", toks)

	parsed := parser.Parser(toks)
	e.logf("parsed: %s\n", parsed)

	transed := parsed
	for _, t := range e.trans {
		transed = t(transed)
	}
	e.logf("transed: %s\n", transed)

	checkEnv := env0.Inherit(e.typeCheck)
	infered := types.TypeCheck(checkEnv, transed)
	e.logf("type: %s\n", infered)

	closure := compile.Compile(e.runtime, transed)
	e.logf("compiled: %v\n", closure)
	return closure
}

func (e *Expr) makeCompiled(closure compile.Closure, env0 *types.Env) Compiled {
	return func(env1 *val.Env) (val *val.Val, err error) {
		defer func() {
			if r := recover(); r != nil {
				if e.dbg != nil {
					e.logf("EVAL FAIL:\n")
					_, _ = e.dbg.Write(debug.Stack())
				}
				err = fmt.Errorf("%s", r)
			}
		}()

		envCheck(env0, env1)
		rt := env1.Inherit(e.runtime)
		val = closure(rt)
		return
	}
}

// envCheck env0 compile-env, env runtime-env
func envCheck(env0 *types.Env, env *val.Env) {
	env0.ForEach(func(name string, kind *types.Kind) {
		v, ok := env.Get(name)
		util.Assert(ok, "undefined %s", name)
		util.Assert(types.Equals(kind, v.Kind),
			"type mismatched, expect `%s` actual `%s`", kind, v.Kind)
	})
}