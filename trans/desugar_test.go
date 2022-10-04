package trans

import (
	. "github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/parser"
	"reflect"
	"testing"
)

func TestDesugar(t *testing.T) {
	tests := []struct {
		expr     string
		expected *Expr
	}{
		{
			expr: "a.b.c(1,2)",
			expected: Call(
				Ident("c"),
				[]*Expr{
					Member(Ident("a"), Ident("b").Ident()),
					LitNum("1"),
					LitNum("2"),
				},
			),
		},
		{
			expr: `"Hello".len()`,
			expected: Call(
				Ident("len"),
				[]*Expr{
					LitStr("Hello"),
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			ops := oper.BuildIn()
			actual := parser.NewParser(ops).Parse(lexer.NewLexer(ops).Lex(tt.expr))
			expected := tt.expected
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expect `%s` actual `%s`", expected, actual)
			}
		})
	}
}
