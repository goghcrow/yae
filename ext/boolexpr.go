package ext

import (
	"github.com/goghcrow/yae/parser/ast"
	"github.com/goghcrow/yae/parser/pos"
)

type LogicalOper int

const (
	AND LogicalOper = iota
	OR
	NOT
)

func (l LogicalOper) String() string { return [...]string{"AND", "OR", "NOT"}[l] }

type BoolExpr interface {
	expr() ast.Expr
}

type Cond struct {
	Field    string
	Operator string
	Operands []ast.Expr
}

type CondGroup struct {
	LogicalOper
	Conds []BoolExpr
}

func (e Cond) expr() ast.Expr {
	args := make([]ast.Expr, len(e.Operands)+1)
	args[0] = ast.Var(e.Field, pos.Unknown)
	copy(args[1:], e.Operands)
	callee := ast.Var(e.Operator, pos.Unknown)
	return ast.Call(callee, args, pos.UnknownCol, pos.Unknown)
}

func (e CondGroup) expr() ast.Expr {
	args := make([]ast.Expr, len(e.Conds))
	for i, c := range e.Conds {
		args[i] = c.expr()
	}
	callee := ast.Var(e.LogicalOper.String(), pos.Unknown)
	return ast.Call(callee, args, pos.UnknownCol, pos.Unknown)
}
