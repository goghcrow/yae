package trans

import (
	"github.com/goghcrow/yae/ast"
)

type Translate func(expr *ast.Expr) *ast.Expr
