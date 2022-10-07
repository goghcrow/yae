package compiler

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/val"
)

type Compiler func(*ast.Expr, *val.Env) Closure

type Closure func(env *val.Env) *val.Val

func (c Closure) String() string { return "Closure" }
