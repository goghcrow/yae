package parser

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
)

func NewParser(ops []oper.Operator) *parser {
	return &parser{grammar: newGrammar(ops)}
}

func (p *parser) Parse(toks []*token.Token) *ast.Expr {
	p.idx = 0
	p.toks = toks
	expr := p.expr(0)
	util.Assert(p.peek() == lexer.EOF, "multi expr")
	return expr
}

type parser struct {
	grammar
	toks []*token.Token
	idx  int
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

func (p *parser) mustEat(typ token.Type) *token.Token {
	t := p.eat()
	util.Assert(t.Type == typ, "syntax error: %s", t)
	return t
}

func (p *parser) tryEat(typ token.Type) *token.Token {
	if p.peek().Type == typ {
		return p.eat()
	} else {
		return nil
	}
}

func (p *parser) tryParse(f func(p *parser) *ast.Expr) (expr *ast.Expr) {
	marked := p.idx
	defer func() {
		if r := recover(); r != nil {
			p.idx = marked
			expr = nil
		}
	}()
	return f(p)
}

func (p *parser) any(fs ...func(p *parser) *ast.Expr) (expr *ast.Expr) {
	for _, f := range fs {
		n := p.tryParse(f)
		if n != nil {
			return n
		}
	}
	p.syntaxAssert(false)
	return nil
}

func (p *parser) begin() *ast.Expr {
	var exprs []*ast.Expr
	for p.peek() == lexer.EOF {
		exprs = append(exprs, p.expr(0))
	}
	return ast.Begin(exprs)
}

func (p *parser) expr(rbp oper.BP) *ast.Expr {
	t := p.eat()
	prefix := p.prefixNud(t)
	left := prefix.nud(p, prefix.BP, t)
	infix := p.parserInfix(left, rbp)
	return p.check(infix)
}

func (p *parser) unary(bp oper.BP, t *token.Token) *ast.Expr {
	return ast.Unary(t.Lexeme, p.expr(bp), true)
}

func (p *parser) unaryP(bp oper.BP, lhs *ast.Expr, t *token.Token) *ast.Expr {
	return ast.Unary(t.Lexeme, lhs, false)
}

func (p *parser) binaryL(bp oper.BP, lhs *ast.Expr, t *token.Token) *ast.Expr {
	rhs := p.expr(bp)
	return ast.Binary(t.Lexeme, oper.INFIX_L, lhs, rhs)
}

func (p *parser) binaryR(bp oper.BP, lhs *ast.Expr, t *token.Token) *ast.Expr {
	rhs := p.expr(bp - 1)
	return ast.Binary(t.Lexeme, oper.INFIX_L, lhs, rhs)
}

func (p *parser) binaryN(bp oper.BP, lhs *ast.Expr, t *token.Token) *ast.Expr {
	rhs := p.expr(bp) // 这里是否-1无所谓, 之后还得过 infixnCheck
	return ast.Binary(t.Lexeme, oper.INFIX_L, lhs, rhs)
}

func (p *parser) parserInfix(left *ast.Expr, rbp oper.BP) *ast.Expr {
	// 判断下一个 tok 是否要绑定 left ( 优先级 > left)
	for p.infixLbp(p.peek()) > rbp {
		t := p.eat()
		infix := p.infixLed(t)
		left = infix.led(p, infix.BP, left, t)
	}
	return left
}

func (p *parser) check(expr *ast.Expr) *ast.Expr {
	if expr.Type == ast.BINARY {
		bin := expr.Binary()
		opName := bin.Name
		if bin.Fixity == oper.INFIX_N {
			if bin.LHS.Type == ast.BINARY {
				util.Assert(bin.LHS.Binary().Name != opName, "non-infix")
			}
			if bin.RHS.Type == ast.BINARY {
				util.Assert(bin.RHS.Binary().Name != opName, "non-infix")
			}
		}
	}
	return expr
}

func (p *parser) syntaxAssert(cond bool) {
	util.Assert(cond, "syntax error: %s", p.toks[p.idx:])
}
