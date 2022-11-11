package test

import (
	"github.com/goghcrow/yae/parser/ast"
	"reflect"
	"testing"
)

type inspector func(ast.Expr) bool

func (f inspector) Visit(expr ast.Expr) ast.Visitor {
	if f(expr) {
		return f
	}
	return nil
}

func Inspect(expr ast.Expr, f func(ast.Expr) bool) {
	ast.Walk(inspector(f), expr)
}

func TestWalk(t *testing.T) {
	expr := parse0("(s.length() + 3) > 5")
	Inspect(expr, func(expr ast.Expr) bool {
		if expr != nil {
			t.Log(reflect.TypeOf(expr))
		}
		t.Log(expr)
		switch expr.(type) {
		case *ast.StrExpr, *ast.NumExpr, *ast.TimeExpr, *ast.BoolExpr, *ast.IdentExpr:
			return false
		}
		return true
	})
}
