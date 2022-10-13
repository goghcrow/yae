package test

import (
	"github.com/goghcrow/yae/ast"
	"testing"
	"unsafe"
)

func TestOffsetOf(t *testing.T) {
	if unsafe.Offsetof(ast.LiteralExpr{}.Expr) != 0 {
		t.Failed()
	}
	if unsafe.Offsetof(ast.LiteralExpr{}.Type) != 0 {
		t.Failed()
	}
}
