package test

import (
	"encoding/json"

	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/compiler"
	"github.com/goghcrow/yae/fun"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/parser"
	"github.com/goghcrow/yae/trans"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
)

func parse0(s string, ops ...oper.Operator) ast.Expr {
	ops = append(oper.BuiltIn(), ops...)
	toks := lexer.NewLexer(ops).Lex(s)
	return parser.NewParser(ops).Parse(toks)
}

func parse(s string, ops ...oper.Operator) ast.Expr {
	return trans.Desugar(parse0(s, ops...))
}

var typecheckEnv = types.NewEnv()
var compileEnv = val.NewEnv()

func init() {
	initEnv(typecheckEnv, compileEnv)
}

func eval(s string, compile compiler.Compiler, typedEnv *types.Env, compileEvalEnv *val.Env) *val.Val {
	toks := lexer.NewLexer(oper.BuiltIn()).Lex(s)
	term := parser.NewParser(oper.BuiltIn()).Parse(toks)
	term = trans.Desugar(term)

	_ = types.Check(term, typedEnv.Inherit(typecheckEnv))
	valuedEnv := compileEvalEnv.Inherit(compileEnv)
	compiled := compile(term, valuedEnv)

	runtimeEnv := val.NewEnv()
	runtimeEnv = runtimeEnv.Inherit(valuedEnv)
	return compiled(runtimeEnv)
}

func infer(s string) *types.Type {
	toks := lexer.NewLexer(oper.BuiltIn()).Lex(s)
	term := parser.NewParser(oper.BuiltIn()).Parse(toks)
	term = trans.Desugar(term)
	return types.Check(term, typecheckEnv)
}

func initEnv(typecheckEnv *types.Env, compileEnv *val.Env) {
	for _, f := range fun.BuiltIn() {
		typecheckEnv.RegisterFun(f.Type)
		compileEnv.RegisterFun(f)
	}
}

func pretty(v interface{}) string {
	s, _ := json.Marshal(v)
	return string(s)
}

func assert(cond bool) {
	if !cond {
		panic(nil)
	}
}
