package parser

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
)

// parser 使用了 Top Down Operator Precedence
// 可以参考道格拉斯的文章 https://www.crockford.com/javascript/tdop/tdop.html

func NewParser(ops []oper.Operator) *parser {
	return &parser{
		grammar: newGrammar(oper.Sort(ops)),
	}
}

func (p *parser) Parse(toks []*token.Token) *ast.Expr {
	p.idx = 0
	p.toks = toks
	expr := p.expr(0)
	p.mustEat(token.EOF)
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
	p.syntaxAssert(t.Type == typ, "expect %s actual %s", typ, t)
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
	util.Assert(false, "try parse fail")
	return nil
}

// parser bp > rbp 的表达式
func (p *parser) expr(rbp oper.BP) *ast.Expr {
	t := p.eat()
	// tok 必须有 prefix 解析器, 否则一定语法错误
	pre := p.mustPrefix(t)
	left := pre.nud(p, pre.BP, t)
	return p.parserInfix(left, rbp)
}

func (p *parser) parserInfix(left *ast.Expr, rbp oper.BP) *ast.Expr {
	// 判断下一个 tok 是否要绑定 left ( 优先级 > left)
	for p.infixLbp(p.peek()) > rbp {
		t := p.eat()
		inf := p.mustInfix(t)
		left = inf.led(p, inf.BP, left, t)
	}
	return p.infixNCheck(left)
}

func (p *parser) infixNCheck(expr *ast.Expr) *ast.Expr {
	if expr.Type == ast.BINARY {
		bin := expr.Binary()
		opName := bin.Name
		if bin.Fixity == oper.INFIX_N {
			if bin.LHS.Type == ast.BINARY {
				p.syntaxAssert(bin.LHS.Binary().Name != opName, "%s non-infix", opName)
			}
			if bin.RHS.Type == ast.BINARY {
				p.syntaxAssert(bin.RHS.Binary().Name != opName, "%s non-infix", opName)
			}
		}
	}
	return expr
}

func (p *parser) syntaxAssert(cond bool, format string, a ...interface{}) {
	util.Assert(cond, "syntax error: "+format, a...)
}
