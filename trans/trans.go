package trans

import (
	"github.com/goghcrow/yae/ast"
)

type Transform func(expr *ast.Expr) *ast.Expr
