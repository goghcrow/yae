package yae

import (
	"fmt"
	"github.com/goghcrow/yae/closure"
	"github.com/goghcrow/yae/compiler"
	"github.com/goghcrow/yae/conv"
	"github.com/goghcrow/yae/fun"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/parser"
	"github.com/goghcrow/yae/trans"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"github.com/goghcrow/yae/vm"
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
	compiled, err := expr.Compile(input, compileTimeEnv)
	if err != nil {
		return nil, err
	}

	runtimeEnv, err := conv.ValEnvOf(v)
	if err != nil {
		return nil, err
	}
	return compiled(runtimeEnv)
}

type Expr struct {
	typeCheck  *types.Env //类型检查环境
	runtime    *val.Env   //编译期运行时环境
	trans      []trans.Translate
	ops        []oper.Operator
	compiler   compiler.Compiler
	useBuildIn bool
	dbg        io.Writer
}

type Callable func(v interface{}) (*val.Val, error)

func NewExpr() *Expr {
	e := Expr{
		typeCheck: types.NewEnv(),
		runtime:   val.NewEnv(),
		trans:     []trans.Translate{},
		//compiler:  closure.Compile,
		compiler:   vm.Compile,
		useBuildIn: true,
	}

	e.initTrans()
	e.initOps()
	e.initFuns()

	return &e
}

func (e *Expr) EnableDebug(out io.Writer) *Expr {
	e.dbg = out
	return e
}

func (e *Expr) UseBytecodeCompiler() *Expr {
	e.compiler = vm.Compile
	return e
}

func (e *Expr) UseClosureCompiler() *Expr {
	e.compiler = closure.Compile
	return e
}

func (e *Expr) UseBuildIn(flag bool) *Expr {
	e.useBuildIn = flag
	return e
}

func (e *Expr) RegisterOperator(ops ...oper.Operator) {
	e.ops = append(e.ops, ops...)
}

func (e *Expr) RegisterTranslator(trans ...trans.Translate) {
	e.trans = append(e.trans, trans...)
}

func (e *Expr) RegisterFun(vs ...*val.Val) {
	for _, v := range vs {
		e.typeCheck.RegisterFun(v.Kind)
		e.runtime.RegisterFun(v)
	}
}

func (e *Expr) Compile(expr string, v interface{}) (c Callable, err error) {
	env0, ok := v.(*types.Env)
	if !ok {
		env0, err = conv.TypeEnvOf(v)
		if err != nil {
			return nil, err
		}
	}
	defer e.backStrace("compile", &err)
	compiled := e.steps(expr, env0)
	c = e.makeCallable(compiled, env0)
	return
}

func (e *Expr) initTrans() {
	e.RegisterTranslator(trans.Desugar)
}

func (e *Expr) initOps() {
	if e.useBuildIn {
		e.ops = oper.BuildIn()
	}
}

func (e *Expr) initFuns() {
	if e.useBuildIn {
		e.RegisterFun(fun.BuildIn()...)
	}
}

func (e *Expr) steps(expr string, env0 *types.Env) compiler.Closure {
	e.logf("expr: %s\n", expr)

	toks := lexer.NewLexer(e.ops).Lex(expr)
	e.logf("lexed: %s\n", toks)

	parsed := parser.NewParser(e.ops).Parse(toks)
	e.logf("parsed: %s\n", parsed)

	transed := parsed
	for _, t := range e.trans {
		transed = t(transed)
	}
	e.logf("transed: %s\n", transed)

	checkEnv := env0.Inherit(e.typeCheck)
	infered := types.Check(transed, checkEnv)
	e.logf("type: %s\n", infered)

	compiled := e.compiler(transed, e.runtime)
	e.logf("compiled: %v\n", compiled)

	return compiled
}

func (e *Expr) makeCallable(closure compiler.Closure, env0 *types.Env) Callable {
	return func(v interface{}) (vl *val.Val, err error) {
		env1, ok := v.(*val.Env)
		if !ok {
			env1, err = conv.ValEnvOf(v)
			if err != nil {
				return nil, err
			}
		}
		defer e.backStrace("eval", &err)
		envCheck(env0, env1)
		rt := env1.Inherit(e.runtime)
		vl = closure(rt)
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

func (e *Expr) backStrace(scene string, err *error) {
	if r := recover(); r != nil {
		if e.dbg != nil {
			e.logf("%s error:\n", scene)
			_, _ = e.dbg.Write(debug.Stack())
		}
		*err = fmt.Errorf("%v", r)
	}
}

func (e *Expr) logf(format string, a ...interface{}) {
	if e.dbg != nil {
		_, _ = fmt.Fprintf(e.dbg, format, a...)
	}
}
