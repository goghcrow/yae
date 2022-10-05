package example

import (
	"fmt"
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/parser"
	"github.com/goghcrow/yae/trans"
	"net/url"
	"os/exec"
	"strings"
	"testing"
)

func parse(s string, ops ...oper.Operator) *ast.Expr {
	ops = append(oper.BuildIn(), ops...)
	toks := lexer.NewLexer(ops).Lex(s)
	return parser.NewParser(ops).Parse(toks)
}

// 可视化 原始语法树 vs 解糖之后的语法树
func TestVisual(t *testing.T) {
	// input := "[(1+2) * -3 / 4, 42, {id:42,name:\"晓\"}.id, [1:2,3:f(1,2,3)][3]][0]"
	input := "(1 + 2) ^ (3 % 4) * 5 - 42 / 100"
	expr := parse(input)
	dot := "sub" + ast.Dot(expr, "cluster1")

	desugarExpr := trans.Desugar(expr)
	desugarDot := "sub" + ast.Dot(desugarExpr, "cluster2")

	buf := &strings.Builder{}
	buf.WriteString("graph \"\" {\n")
	buf.WriteString(fmt.Sprintf("label=%q\n", input))
	buf.WriteString(dot)
	buf.WriteString("\n")
	buf.WriteString(desugarDot)
	buf.WriteString("\n")
	buf.WriteString("}")

	s := buf.String()
	t.Log(s)

	cmd := exec.Command("open", "https://dreampuf.github.io/GraphvizOnline/#"+url.PathEscape(s))
	_, _ = cmd.Output()
}