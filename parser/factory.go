package parser

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/token"
)

func newGrammar() grammar {
	var g = grammar{}
	g.prefix(token.NAME, func(p *parser, t token.Token) *ast.Expr { return ast.Ident(t.Lexeme) })

	g.prefix(token.NULL, literal(ast.LIT_NULL))
	g.prefix(token.TRUE, literal(ast.LIT_TRUE))
	g.prefix(token.FALSE, literal(ast.LIT_FALSE))
	g.prefix(token.NUM, literal(ast.LIT_NUM))
	g.prefix(token.STR, literal(ast.LIT_STR))

	g.prefix(token.LEFT_BRACKET, parseListMap)
	g.prefix(token.LEFT_BRACE, parseObj)
	g.prefix(token.LEFT_PAREN, parseGroup)

	g.prefix(token.PLUS, unary)  // token.UNARY_PLUS
	g.prefix(token.MINUS, unary) // token.UNARY_MINUS

	g.infixLeft(token.MUL, binaryL)
	g.infixLeft(token.DIV, binaryL)
	g.infixLeft(token.MOD, binaryL)
	g.infixLeft(token.PLUS, binaryL)
	g.infixLeft(token.MINUS, binaryL)
	g.infixLeft(token.EXP, binaryR)

	g.infixLeft(token.GT, binaryL)
	g.infixLeft(token.LT, binaryL)
	g.infixLeft(token.GE, binaryL)
	g.infixLeft(token.LE, binaryL)
	g.infixLeft(token.EQ, binaryL)
	g.infixLeft(token.NE, binaryL)

	g.prefix(token.LOGIC_NOT, unary)
	g.infixLeft(token.LOGIC_AND, binaryL)
	g.infixLeft(token.LOGIC_OR, binaryL)
	g.infixRight(token.QUESTION, parseQuestion)
	g.infixLeft(token.DOT, parseDot)

	g.prefix(token.IF, parseIf)

	g.infixLeft(token.LEFT_PAREN, parseCall)
	g.infixLeft(token.LEFT_BRACKET, parseSubscript)
	return g
}

func literal(typ ast.LitType) nud {
	return func(p *parser, t token.Token) *ast.Expr { return ast.Literal(typ, t.Lexeme) }
}
func binaryN(p *parser, lhs *ast.Expr, t token.Token) *ast.Expr { return p.binaryN(lhs, t) }
func binaryL(p *parser, lhs *ast.Expr, t token.Token) *ast.Expr { return p.binaryL(lhs, t) }
func binaryR(p *parser, lhs *ast.Expr, t token.Token) *ast.Expr { return p.binaryR(lhs, t) }
func unary(p *parser, t token.Token) *ast.Expr                  { return p.unary(t) }
func unaryP(p *parser, lhs *ast.Expr, t token.Token) *ast.Expr  { return p.unaryP(lhs, t) }

func parseListMap(p *parser, t token.Token) *ast.Expr {
	if p.tryEat(token.COLON) != nil {
		p.mustEat(token.RIGHT_BRACKET)
		return ast.Map([]ast.Pair{})
	}
	return p.any(parseList, parseMap)
}

func parseList(p *parser) *ast.Expr {
	elems := make([]*ast.Expr, 0)
	for {
		if p.peek().Type == token.RIGHT_BRACKET {
			break
		}
		elems = append(elems, p.expr(0))
		if p.tryEat(token.COMMA) == nil {
			break
		}
	}
	p.mustEat(token.RIGHT_BRACKET)
	return ast.List(elems)
}

func parseMap(p *parser) *ast.Expr {
	pairs := make([]ast.Pair, 0)
	for {
		if p.peek().Type == token.RIGHT_BRACKET {
			break
		}
		k := p.expr(0)
		p.mustEat(token.COLON)
		v := p.expr(0)
		pairs = append(pairs, ast.Pair{Key: k, Val: v})
		if p.tryEat(token.COMMA) == nil {
			break
		}
	}
	p.mustEat(token.RIGHT_BRACKET)
	return ast.Map(pairs)
}

func parseObj(p *parser, t token.Token) *ast.Expr {
	fs := make(map[string]*ast.Expr, 0)
	for {
		if p.peek().Type == token.RIGHT_BRACE {
			break
		}
		n := p.mustEat(token.NAME)
		p.mustEat(token.COLON)
		v := p.expr(0)
		fs[n.Lexeme] = v
		if p.tryEat(token.COMMA) == nil {
			break
		}
	}
	p.mustEat(token.RIGHT_BRACE)
	return ast.Obj(fs)
}

func parseGroup(p *parser, t token.Token) *ast.Expr {
	expr := p.expr(0)
	p.mustEat(token.RIGHT_PAREN)
	return expr
}

// if expr then expr else xxx [end]
func parseIf(p *parser, iff token.Token) *ast.Expr {
	cond := p.expr(0)
	p.mustEat(token.THEN)
	then := p.expr(0)
	p.mustEat(token.ELSE)
	els := p.expr(0)
	p.tryEat(token.END)
	return ast.If(cond, then, els)
}

func parseQuestion(p *parser, l *ast.Expr, t token.Token) *ast.Expr {
	m := p.expr(token.QUESTION.Bp())
	p.mustEat(token.COLON)
	r := p.expr(token.QUESTION.Bp() - 1)
	return ast.Tenary(token.QUESTION, l, m, r)
}

func parseCall(p *parser, callee *ast.Expr, t token.Token) *ast.Expr {
	args := make([]*ast.Expr, 0)
	rp := p.tryEat(token.RIGHT_PAREN)
	if rp == nil {
		for {
			args = append(args, p.expr(0))
			if p.tryEat(token.COMMA) == nil {
				break
			}
		}
		p.mustEat(token.RIGHT_PAREN)
	}
	return ast.Call(callee, args)
}

func parseDot(p *parser, obj *ast.Expr, t token.Token) *ast.Expr {
	name := p.eat()
	// 放开限制则可以写 1.+(1), 1可以看成对象
	//util.Assert(name.Type == token.NAME || name.Type == token.TRUE || name.Type == token.FALSE,
	//	"syntax error: %s", name.Lexeme)
	expr := ast.Member(obj, ast.Ident(name.Lexeme).Ident())
	lp := p.tryEat(token.LEFT_PAREN)
	if lp == nil {
		return expr
	}
	return parseCall(p, expr, t)
}

func parseSubscript(p *parser, list *ast.Expr, t token.Token) *ast.Expr {
	expr := p.expr(0)
	p.mustEat(token.RIGHT_BRACKET)
	return ast.Subscript(list, expr)
}
