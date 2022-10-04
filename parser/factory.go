package parser

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
)

func newGrammar(ops []oper.Operator) grammar {
	g := grammar{
		prefixs: map[token.Type]prefix{},
		infixs:  map[token.Type]infix{},
	}

	g.prefix(token.NAME, oper.BP_NONE, ident)

	g.prefix(token.NULL, oper.BP_NONE, literal(ast.LIT_NULL))
	g.prefix(token.TRUE, oper.BP_NONE, literal(ast.LIT_TRUE))
	g.prefix(token.FALSE, oper.BP_NONE, literal(ast.LIT_FALSE))
	g.prefix(token.NUM, oper.BP_NONE, literal(ast.LIT_NUM))
	g.prefix(token.STR, oper.BP_NONE, literal(ast.LIT_STR))
	g.prefix(token.TIME, oper.BP_NONE, literal(ast.LIT_TIME))

	g.prefix(token.LEFT_BRACKET, oper.BP_NONE, parseListMap)
	g.prefix(token.LEFT_BRACE, oper.BP_NONE, parseObj)
	g.prefix(token.LEFT_PAREN, oper.BP_NONE, parseGroup)

	g.prefix(token.IF, oper.BP_NONE, parseIf) // 如果 if 是普通函数, 这里可以干掉

	for _, op := range ops {
		switch op.Fixity {
		case oper.PREFIX:
			g.prefix(op.Type, op.BP, unaryPrefix)
		case oper.INFIX_N:
			g.infix(op.Type, op.BP, binaryN)
		case oper.INFIX_L:
			g.infix(op.Type, op.BP, binaryL)
		case oper.INFIX_R:
			g.infix(op.Type, op.BP, binaryR)
		case oper.POSTFIX:
			g.postfix(op.Type, op.BP, unaryPostfix)
		}
	}

	// 放在自定义操作符后面, 防止 ? . 被覆盖
	g.infixRight(token.QUESTION, oper.BP_COND, parseQuestion)
	g.infixLeft(token.DOT, oper.BP_MEMBER, parseDot)

	g.infixLeft(token.LEFT_PAREN, oper.BP_CALL, parseCall)
	g.infixLeft(token.LEFT_BRACKET, oper.BP_MEMBER, parseSubscript)
	return g
}

func ident(p *parser, bp oper.BP, t *token.Token) *ast.Expr {
	return ast.Ident(t.Lexeme)
}

func literal(typ ast.LitType) nud {
	return func(p *parser, bp oper.BP, t *token.Token) *ast.Expr { return ast.Literal(typ, t.Lexeme) }
}

func binaryL(p *parser, bp oper.BP, lhs *ast.Expr, t *token.Token) *ast.Expr {
	rhs := p.expr(bp)
	return ast.Binary(t.Lexeme, oper.INFIX_L, lhs, rhs)
}

func binaryR(p *parser, bp oper.BP, lhs *ast.Expr, t *token.Token) *ast.Expr {
	rhs := p.expr(bp - 1)
	return ast.Binary(t.Lexeme, oper.INFIX_R, lhs, rhs)
}

func binaryN(p *parser, bp oper.BP, lhs *ast.Expr, t *token.Token) *ast.Expr {
	rhs := p.expr(bp) // 这里是否-1无所谓, 之后会检查
	return ast.Binary(t.Lexeme, oper.INFIX_N, lhs, rhs)
}

func unaryPrefix(p *parser, bp oper.BP, t *token.Token) *ast.Expr {
	expr := p.expr(bp)
	return ast.Unary(t.Lexeme, expr, true)
}

func unaryPostfix(p *parser, bp oper.BP, lhs *ast.Expr, t *token.Token) *ast.Expr {
	return ast.Unary(t.Lexeme, lhs, false)
}

func parseListMap(p *parser, bp oper.BP, t *token.Token) *ast.Expr {
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
		el := p.expr(0)
		elems = append(elems, el)
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

func parseObj(p *parser, bp oper.BP, t *token.Token) *ast.Expr {
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

func parseGroup(p *parser, bp oper.BP, t *token.Token) *ast.Expr {
	expr := p.expr(0)
	p.mustEat(token.RIGHT_PAREN)
	return expr
}

func parseQuestion(p *parser, bp oper.BP, l *ast.Expr, t *token.Token) *ast.Expr {
	m := p.expr(0)
	p.mustEat(token.COLON)
	r := p.expr(bp - 1)
	return ast.Tenary(token.QUESTION, l, m, r)
}

func parseCall(p *parser, bp oper.BP, callee *ast.Expr, t *token.Token) *ast.Expr {
	args := make([]*ast.Expr, 0)
	rp := p.tryEat(token.RIGHT_PAREN)
	if rp == nil {
		for {
			arg := p.expr(0)
			args = append(args, arg)
			if p.tryEat(token.COMMA) == nil {
				break
			}
		}
		p.mustEat(token.RIGHT_PAREN)
	}
	return ast.Call(callee, args)
}

func parseDot(p *parser, bp oper.BP, obj *ast.Expr, t *token.Token) *ast.Expr {
	name := p.eat()
	// 放开限制则可以写 1. +(1), 1可以看成对象, .和+必须有空格是因为否则会匹配自定义操作符
	//util.Assert(name.Type == token.NAME || name.Type == token.TRUE || name.Type == token.FALSE,
	//	"syntax error: %s", name.Lexeme)
	expr := ast.Member(obj, ast.Ident(name.Lexeme).Ident())
	lp := p.tryEat(token.LEFT_PAREN)
	if lp == nil {
		return expr
	}
	return parseCall(p, bp, expr, t)
}

func parseSubscript(p *parser, bp oper.BP, list *ast.Expr, t *token.Token) *ast.Expr {
	expr := p.expr(0)
	p.mustEat(token.RIGHT_BRACKET)
	return ast.Subscript(list, expr)
}

// if expr then expr else xxx [end]
func parseIf(p *parser, bp oper.BP, iff *token.Token) *ast.Expr {
	cond := p.expr(0)
	p.mustEat(token.THEN)
	then := p.expr(0)
	p.mustEat(token.ELSE)
	els := p.expr(0)
	p.tryEat(token.END)
	return ast.If(cond, then, els)
}
