package facade

import (
	"errors"
	"fmt"
	"github.com/goghcrow/yae/check"
	"github.com/goghcrow/yae/compile"
	"github.com/goghcrow/yae/env"
	"github.com/goghcrow/yae/env0"
	lex "github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/parser"
	"github.com/goghcrow/yae/val"
	"runtime/debug"
)

type Compiled func(env1 *env.Env) (*val.Val, error)

func (e *Expr) Compile(expr string, env0 *env0.Env) (r Compiled, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e.dbg != nil {
				e.logf("COMPILE FAIL:\n")
				_, _ = e.dbg.Write(debug.Stack())
			}
			err = errors.New(fmt.Sprintf("%s", r))
		}
	}()
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
	infered := check.TypeCheck(checkEnv, transed)
	e.logf("type: %s\n", infered)
	closure := compile.Compile(e.runtime, transed)
	e.logf("compiled: %v\n", closure)

	r = func(env1 *env.Env) (val *val.Val, err error) {
		defer func() {
			if r := recover(); r != nil {
				if e.dbg != nil {
					e.logf("EVAL FAIL:\n")
					_, _ = e.dbg.Write(debug.Stack())
				}
				err = errors.New(fmt.Sprintf("%s", r))
			}
		}()

		check.EnvCheck(env0, env1)
		rt := env1.Inherit(e.runtime)
		val = closure(rt)
		return
	}
	return
}
