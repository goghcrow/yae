package parser

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/trans"
	"testing"
)

func lex(input string) []*token.Token {
	return lexer.NewLexer(oper.BuildIn()).Lex(input)
}

func parse(toks []*token.Token) *ast.Expr {
	return NewParser(oper.BuildIn()).Parse(toks)
}

func TestParser(t *testing.T) {
	{
		toks := lex(`-42 == 1`)
		t.Log(toks)
		ast := trans.Desugar(parse(toks))
		t.Log(ast)
	}

	//{
	//	// todo
	//	// 遇到一个比他小的前缀操作符应该结合不了!!!
	//	// unary 构造器里头 实际上是 用 binary 的 lbp 去取表达式, 而不是实际的 0
	//	//
	//	// return ast.Unary(t.Type, p.expr(t.Type.BP), true), 这里就应该直接写 0 ???, 啥前缀都能
	//
	//	// 小于 binary - 的优先级的操作符，会发生结合不上得问题
	//
	//	def := append(oper.BuildIn(), &token.Type{
	//		Name: "#",
	//		//BP:     token.BP_TERM + 1, // -(#(1, 2))
	//		//BP:     token.BP_TERM, // #(-(1), 2)
	//		//BP:     token.BP_TERM - 1, // #(-(1), 2)
	//		BP:     1,
	//		Fixity: token.INFIX_L,
	//	})
	//
	//	// `-1  2`
	//	toks := lexer.NewLexer(def).Lex(`-1 # 2`)
	//	t.Log(toks)
	//
	//	ast := trans.Desugar(NewParser(def).Parse(toks))
	//	t.Log(ast)
	//}

	//{
	//	toks := lexer.Lex(`1 - 2 + 3 * 4`)
	//	t.Log(toks)
	//	ast := trans.Desugar(Parse(toks))
	//	t.Log(ast)
	//}
	//{
	//	toks := lexer.Lex(`+1`)
	//	t.Log(toks)
	//	ast := trans.Desugar(Parse(toks))
	//	t.Log(ast)
	//}
	//{
	//	toks := lexer.Lex(`1 + -1`)
	//	t.Log(toks)
	//	ast := trans.Desugar(Parse(toks))
	//	t.Log(ast)
	//}
	//{
	//	toks := lexer.Lex(`1 - -1`)
	//	t.Log(toks)
	//	ast := trans.Desugar(Parse(toks))
	//	t.Log(ast)
	//}
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
	return parse(lex(s))
}
