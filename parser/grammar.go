package parser

import (
	"github.com/goghcrow/yae/parser/ast"
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/token"
	"github.com/goghcrow/yae/util"
)

type nud func(*parser, oper.BP, *token.Token) ast.Expr
type led func(*parser, oper.BP, ast.Expr, *token.Token) ast.Expr

// 处理字面量、变量、前缀操作符
type prefix struct {
	oper.BP
	nud
}

// 处理中缀、后缀操作符
type infix struct {
	oper.BP
	led
}

// 如果不支持自定操作符, 则 tokenType 可以定义成 int enum
// prefixs & infixs 则可以定义成 tokenType 为下标的数组

type grammar struct {
	prefixs map[token.Kind]prefix
	infixs  map[token.Kind]infix
}

// 前缀操作符
func (g *grammar) prefix(k token.Kind, bp oper.BP, f nud) {
	g.prefixs[k] = prefix{bp, f}
}

// 不结合中缀操作符
func (g *grammar) infix(k token.Kind, bp oper.BP, f led) {
	g.infixs[k] = infix{bp, f}
}

// 右结合中缀操作符
func (g *grammar) infixRight(k token.Kind, bp oper.BP, f led) {
	g.infix(k, bp, f)
}

// 左结合中缀操作符
func (g *grammar) infixLeft(k token.Kind, bp oper.BP, f led) {
	g.infix(k, bp, f)
}

// 后缀操作符（可以看成中缀操作符木有右边操作数）
func (g *grammar) postfix(k token.Kind, bp oper.BP, f led) {
	g.infix(k, bp, f)
}

// left binding powers
func (p *grammar) infixLbp(t *token.Token) oper.BP {
	i, ok := p.infixs[t.Kind]
	if ok {
		return i.BP
	} else {
		return 0
	}
}

func (g *grammar) mustPrefix(t *token.Token) prefix {
	p, ok := g.prefixs[t.Kind]
	util.Assert(ok, "syntax error in %s: %s", t.Loc, t)
	return p
}

func (g *grammar) mustInfix(t *token.Token) infix {
	i, ok := g.infixs[t.Kind]
	util.Assert(ok, "syntax error in %s: %s", t.Loc, t)
	return i
}
