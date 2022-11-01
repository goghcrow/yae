package trans

import (
	"github.com/goghcrow/yae/parser/ast"
)

type Translate func(expr ast.Expr) ast.Expr
