package parser

import (
	"github.com/goghcrow/yae/parser/ast"
	"github.com/goghcrow/yae/parser/lexer"
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/pos"
	"github.com/goghcrow/yae/parser/token"
	"github.com/goghcrow/yae/util"
)

// parser 使用了 Top Down Operator Precedence
// 可以参考道格拉斯的文章 https://www.crockford.com/javascript/tdop/tdop.html
// 以及论文 https://tdop.github.io/
// 以及另一个名称 Pratt Parsers,
// https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/

type parser struct {
	grammar
	toks []*token.Token
	idx  int
}

func NewParser(ops []oper.Operator) *parser {
	return &parser{
		grammar: newGrammar(ops),
	}
}

func (p *parser) Parse(toks []*token.Token) ast.Expr {
	p.idx = 0
	p.toks = toks

	expr := p.expr(0)
	p.expectEOF()
	return expr
}

func (p *parser) expectEOF() {
	p.mustEat(token.EOF)
}

func (p *parser) expectSingleExpr() {
	// 可以支持用 ;; 分隔 topLevel 表达式
	// 如果用 \n 分隔, 需要在 parser 各处 tryEatLines
}

func (p *parser) peek() *token.Token {
	if p.idx >= len(p.toks) {
		return lexer.EOF
	}
	return p.toks[p.idx]
}

func (p *parser) eat() *token.Token {
	if p.idx >= len(p.toks) {
		return lexer.EOF
	}
	t := p.toks[p.idx]
	p.idx++
	return t
}

func (p *parser) mustEat(k token.Kind) *token.Token {
	t := p.eat()
	p.syntaxAssert(t.Pos, t.Kind == k, "expect `%s` actual `%s`", k, t)
	return t
}

func (p *parser) tryEat(k token.Kind) *token.Token {
	if p.peek().Kind == k {
		return p.eat()
	} else {
		return nil
	}
}

func (p *parser) tryParse(f func(p *parser) ast.Expr) (expr ast.Expr) {
	marked := p.idx
	defer func() {
		if r := recover(); r != nil {
			p.idx = marked
			expr = nil
		}
	}()
	return f(p)
}

func (p *parser) any(expect string, fs ...func(p *parser) ast.Expr) (expr ast.Expr) {
	for _, f := range fs {
		n := p.tryParse(f)
		if n != nil {
			return n
		}
	}
	p.syntaxAssert(p.peek().Pos, false, "expect `%s`", expect)
	return nil
}

// parser bp > rbp 的表达式
func (p *parser) expr(rbp oper.BP) ast.Expr {
	t := p.eat()
	// tok 必须有 prefix 解析器, 否则一定语法错误
	pre := p.mustPrefix(t)
	left := pre.nud(p, pre.BP, t)
	return p.parseInfix(left, rbp)
}

func (p *parser) parseInfix(left ast.Expr, rbp oper.BP) ast.Expr {
	// 判断下一个 tok 是否要绑定 left ( 优先级 > left)
	for p.infixLbp(p.peek()) > rbp {
		t := p.eat()
		inf := p.mustInfix(t)
		left = inf.led(p, inf.BP, left, t)
	}
	return p.infixNCheck(left)
}

func (p *parser) infixNCheck(expr ast.Expr) ast.Expr {
	if bin, ok := expr.(*ast.BinaryExpr); ok {
		opName := bin.Name
		if bin.Fixity == oper.INFIX_N {
			if lhs, ok := bin.LHS.(*ast.BinaryExpr); ok {
				p.syntaxAssert(lhs.IdentExpr.Pos, lhs.Name != opName, "%s non-infix", opName)
			}
			if rhs, ok := bin.RHS.(*ast.BinaryExpr); ok {
				p.syntaxAssert(rhs.IdentExpr.Pos, rhs.Name != opName, "%s non-infix", opName)
			}
		}
	}
	return expr
}

func (p *parser) syntaxAssert(l pos.Pos, cond bool, format string, a ...interface{}) {
	util.Assert(cond, "syntax error in "+l.String()+": "+format, a...)
}
