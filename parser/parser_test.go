package parser

import (
	"github.com/goghcrow/yae/ast"
	lex "github.com/goghcrow/yae/lexer"
	"testing"
)

func TestParser(t *testing.T) {
	Parser(lex.Lex(`'now'`))
}

func TestSyntaxError(t *testing.T) {
	for _, expr := range []string{
		`"Hello" + `,
	} {
		ast := syntaxError(expr)
		if ast != nil {
			t.Errorf("expected syntax error get `%s`", ast)
		}
	}
}

func syntaxError(s string) (e *ast.Expr) {
	defer func() {
		if r := recover(); r != nil {
			e = nil
		}
	}()
	return Parser(lex.Lex(s))
}
