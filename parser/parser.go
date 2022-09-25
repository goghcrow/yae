package parser

import (
	"github.com/goghcrow/yae/ast"
	lex "github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
)

var g = newGrammar()

func Parser(toks []*token.Token) *ast.Expr {
	p := parser{toks: toks, grammar: g}
	expr := p.expr(0)
	util.Assert(p.peek() == lex.EOF, "multi expr")
	return expr
}

// 处理字面量、变量、前缀操作符
type nud func(*parser, token.Token) *ast.Expr

// 处理中缀、后缀操作符
type led func(*parser, *ast.Expr, token.Token) *ast.Expr

type parser struct {
	grammar
	toks []*token.Token
	idx  int
}

func (p *parser) peek() token.Token {
	if p.idx >= len(p.toks) {
		return lex.EOF
	}
	return *p.toks[p.idx]
}
func (p *parser) eat() token.Token {
	if p.idx >= len(p.toks) {
		return lex.EOF
	}
	t := p.toks[p.idx]
	p.idx++
	return *t
}
func (p *parser) mustEat(typ token.Type) token.Token {
	t := p.eat()
	util.Assert(t.Type == typ, "syntax error: %s", &t)
	return t
}
func (p *parser) tryEat(typ token.Type) *token.Token {
	if p.peek().Type == typ {
		t := p.eat()
		return &t
	} else {
		return nil
	}
}

//func (p *parser) begin() *ast.Expr {
//	var exprs []*ast.Expr
//	for p.peek() == lex.EOF {
//		exprs = append(exprs, p.expr(0))
//	}
//	return ast.Begin(exprs)
//}

func (p *parser) expr(rbp token.BP) *ast.Expr {
	t := p.eat()
	left := p.prefixNud(t)(p, t)
	infix := p.parserInfix(left, rbp)
	return p.check(infix)
}

func (p *parser) unary(t token.Token) *ast.Expr {
	return ast.Unary(t.Type, p.expr(t.Type.Bp()), true)
}

func (p *parser) unaryP(lhs *ast.Expr, t token.Token) *ast.Expr {
	return ast.Unary(t.Type, lhs, false)
}

func (p *parser) binaryL(lhs *ast.Expr, t token.Token) *ast.Expr {
	rhs := p.expr(t.Type.Bp())
	return ast.Binary(t.Type, lhs, rhs)
}

func (p *parser) binaryR(lhs *ast.Expr, t token.Token) *ast.Expr {
	rhs := p.expr(t.Type.Bp() - 1)
	return ast.Binary(t.Type, lhs, rhs)
}

func (p *parser) binaryN(lhs *ast.Expr, t token.Token) *ast.Expr {
	// 这里 binaryL + R 其实无所谓, 之后还得过 infixnCheck
	return p.binaryL(lhs, t)
}

func (p *parser) parserInfix(left *ast.Expr, rbp token.BP) *ast.Expr {
	// 判断下一个 tok 是否要绑定 left ( 优先级 > left)
	for p.infixLbp(p.peek()) > rbp {
		t := p.eat()
		left = p.infixLed(t)(p, left, t)
	}
	return left
}

func (p *parser) prefixNud(tok token.Token) nud {
	n := p.prefixs[tok.Type]
	util.Assert(n != nil, "syntax error: %s", &tok)
	return n
}

func (p *parser) infixLed(tok token.Token) led {
	l := p.infixs[tok.Type]
	util.Assert(l != nil, "syntax error: %s", &tok)
	return l
}

func (p *parser) infixLbp(tok token.Token) token.BP {
	return p.lbps[tok.Type]
}

func (p *parser) check(expr *ast.Expr) *ast.Expr {
	if expr.Type == ast.BINARY {
		bin := expr.Binary()
		op := bin.Type.Name()
		if bin.Fixity() == token.INFIX_N {
			if bin.LHS.Type == ast.BINARY {
				util.Assert(bin.LHS.Binary().Type.Name() != op, "non-infix")
			}
			if bin.RHS.Type == ast.BINARY {
				util.Assert(bin.RHS.Binary().Type.Name() != op, "non-infx")
			}
		}
	}
	return expr
}
