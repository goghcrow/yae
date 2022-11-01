package parser

import (
	"github.com/goghcrow/yae/parser/ast"
	"github.com/goghcrow/yae/parser/loc"
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/token"
)

func newGrammar(ops []oper.Operator) grammar {
	g := grammar{
		prefixs: map[token.Type]prefix{},
		infixs:  map[token.Type]infix{},
	}

	// 这里如果 token.Type 不重复, 顺序无所谓, 重复默认覆盖

	g.prefix(token.NAME, oper.BP_NONE, ident)

	g.prefix(token.TRUE, oper.BP_NONE, parseTrue)
	g.prefix(token.FALSE, oper.BP_NONE, parseFalse)
	g.prefix(token.NUM, oper.BP_NONE, parseNum)
	g.prefix(token.STR, oper.BP_NONE, parseStr)
	g.prefix(token.TIME, oper.BP_NONE, parseTime)

	g.prefix(token.LEFT_BRACKET, oper.BP_NONE, parseListMap)
	g.prefix(token.LEFT_BRACE, oper.BP_NONE, parseObj)
	g.prefix(token.LEFT_PAREN, oper.BP_NONE, parseGroup)

	// if 是普通函数
	// g.prefix(token.IF, oper.BP_NONE, parseIf)

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

func parseTrue(p *parser, bp oper.BP, t *token.Token) ast.Expr  { return ast.True(t.Loc) }
func parseFalse(p *parser, bp oper.BP, t *token.Token) ast.Expr { return ast.False(t.Loc) }
func parseNum(p *parser, bp oper.BP, t *token.Token) ast.Expr   { return ast.Num(t.Lexeme, t.Loc) }
func parseStr(p *parser, bp oper.BP, t *token.Token) ast.Expr   { return ast.Str(t.Lexeme, t.Loc) }
func parseTime(p *parser, bp oper.BP, t *token.Token) ast.Expr  { return ast.Time(t.Lexeme, t.Loc) }

func ident(p *parser, bp oper.BP, t *token.Token) ast.Expr {
	return ast.Var(t.Lexeme, t.Loc)
}

func binaryL(p *parser, bp oper.BP, lhs ast.Expr, t *token.Token) ast.Expr {
	rhs := p.expr(bp)
	loc := lhs.GetLoc().Merge(rhs.GetLoc())
	return ast.Binary(ast.Var(t.Lexeme, t.Loc), oper.INFIX_L, lhs, rhs, loc)
}

func binaryR(p *parser, bp oper.BP, lhs ast.Expr, t *token.Token) ast.Expr {
	rhs := p.expr(bp - 1)
	loc := lhs.GetLoc().Merge(rhs.GetLoc())
	return ast.Binary(ast.Var(t.Lexeme, t.Loc), oper.INFIX_R, lhs, rhs, loc)
}

func binaryN(p *parser, bp oper.BP, lhs ast.Expr, t *token.Token) ast.Expr {
	rhs := p.expr(bp) // 这里是否-1无所谓, 之后会检查
	loc := lhs.GetLoc().Merge(rhs.GetLoc())
	return ast.Binary(ast.Var(t.Lexeme, t.Loc), oper.INFIX_N, lhs, rhs, loc)
}

func unaryPrefix(p *parser, bp oper.BP, t *token.Token) ast.Expr {
	expr := p.expr(bp)
	loc := t.GetLoc().Merge(expr.GetLoc())
	return ast.Unary(ast.Var(t.Lexeme, t.Loc), expr, true, loc)
}

func unaryPostfix(p *parser, bp oper.BP, lhs ast.Expr, t *token.Token) ast.Expr {
	loc := lhs.GetLoc().Merge(t.GetLoc())
	return ast.Unary(ast.Var(t.Lexeme, t.Loc), lhs, false, loc)
}

func parseListMap(p *parser, bp oper.BP, t *token.Token) ast.Expr {
	if p.tryEat(token.COLON) != nil {
		rb := p.mustEat(token.RIGHT_BRACKET)
		loc := t.GetLoc().Merge(rb.GetLoc())
		return ast.Map([]ast.Pair{}, loc)
	}
	return p.any("list or map", parseList(t), parseMap(t))
}

func parseList(t *token.Token) func(p *parser) ast.Expr {
	return func(p *parser) ast.Expr {
		elems := make([]ast.Expr, 0)
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
		rb := p.mustEat(token.RIGHT_BRACKET)
		loc := t.GetLoc().Merge(rb.GetLoc())
		return ast.List(elems, loc)
	}
}

func parseMap(t *token.Token) func(p *parser) ast.Expr {
	return func(p *parser) ast.Expr {
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
		rb := p.mustEat(token.RIGHT_BRACKET)
		loc := t.GetLoc().Merge(rb.GetLoc())
		return ast.Map(pairs, loc)
	}
}

func parseObj(p *parser, bp oper.BP, t *token.Token) ast.Expr {
	fs := make([]ast.Field, 0)
	for {
		if p.peek().Type == token.RIGHT_BRACE {
			break
		}
		n := p.mustEat(token.NAME)
		p.mustEat(token.COLON)
		v := p.expr(0)
		fs = append(fs, ast.Field{Name: n.Lexeme, Val: v})
		if p.tryEat(token.COMMA) == nil {
			break
		}
	}
	rb := p.mustEat(token.RIGHT_BRACE)
	loc := t.GetLoc().Merge(rb.GetLoc())
	return ast.Obj(fs, loc)
}

func parseGroup(p *parser, bp oper.BP, t *token.Token) ast.Expr {
	expr := p.expr(0)
	rp := p.mustEat(token.RIGHT_PAREN)
	loc := t.GetLoc().Merge(rp.GetLoc())
	return ast.Group(expr, loc)
}

func parseQuestion(p *parser, bp oper.BP, l ast.Expr, t *token.Token) ast.Expr {
	m := p.expr(0)
	p.mustEat(token.COLON)
	r := p.expr(bp - 1)
	loc := l.GetLoc().Merge(r.GetLoc())
	return ast.Tenary(ast.Var(t.Lexeme, t.Loc), l, m, r, loc)
}

func parseCall(p *parser, bp oper.BP, callee ast.Expr, t *token.Token) ast.Expr {
	args := make([]ast.Expr, 0)
	rp := p.tryEat(token.RIGHT_PAREN)
	if rp == nil {
		for {
			arg := p.expr(0)
			args = append(args, arg)
			if p.tryEat(token.COMMA) == nil {
				break
			}
		}
		rp = p.mustEat(token.RIGHT_PAREN)
	}
	callLoc := callee.GetLoc().Merge(rp.GetLoc())
	return ast.Call(callee, args, loc.DbgCol(t.Col), callLoc)
}

func parseDot(p *parser, bp oper.BP, obj ast.Expr, t *token.Token) ast.Expr {
	name := p.eat()
	// 放开限制则可以写 1. +(1), 1可以看成对象, .和+必须有空格是因为否则会匹配自定义操作符
	//util.Assert(name.Type == token.NAME || name.Type == token.TRUE || name.Type == token.FALSE,
	//	"syntax error: %s", name.Lexeme)
	selLoc := obj.GetLoc().Merge(name.GetLoc())
	expr := ast.Member(obj, ast.Var(name.Lexeme, name.Loc), loc.DbgCol(t.Col), selLoc)
	lp := p.tryEat(token.LEFT_PAREN)
	if lp == nil {
		return expr
	} else {
		return parseCall(p, bp, expr, lp)
	}
}

func parseSubscript(p *parser, bp oper.BP, list ast.Expr, t *token.Token) ast.Expr {
	expr := p.expr(0)
	rb := p.mustEat(token.RIGHT_BRACKET)
	subsLoc := list.GetLoc().Merge(rb.GetLoc())
	return ast.Subscript(list, expr, loc.DbgCol(t.Col), subsLoc)
}

// if expr then expr else xxx [end]
//func parseIf(p *parser, bp oper.BP, iff *token.Token) ast.Expr {
//	cond := p.expr(0)
//	p.mustEat(token.THEN)
//	then := p.expr(0)
//	p.mustEat(token.ELSE)
//	els := p.expr(0)
//	end := p.tryEat(token.END)
//	var loc loc.Loc
//	if end == nil {
//		loc = iff.GetLoc().Merge(els.GetLoc())
//	} else {
//		loc = iff.GetLoc().Merge(end.GetLoc())
//	}
//	return ast.If(cond, then, els, loc)
//}
