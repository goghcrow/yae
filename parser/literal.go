package parser

import (
	"errors"
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/timelib"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
	"strconv"
	"strings"
)

func parseTrue(p *parser, bp oper.BP, t *token.Token) ast.Expr  { return ast.True() }
func parseFalse(p *parser, bp oper.BP, t *token.Token) ast.Expr { return ast.False() }
func parseNum(p *parser, bp oper.BP, t *token.Token) ast.Expr {
	f, err := parseNum0(t.Lexeme)
	util.Assert(err == nil, "invalid num literal %s", t.Lexeme)
	return ast.Num(t.Lexeme, f)
}
func parseStr(p *parser, bp oper.BP, t *token.Token) ast.Expr {
	unquoted, err := strconv.Unquote(t.Lexeme)
	util.Assert(err == nil, "invalid string literal: %s", t.Lexeme)
	return ast.Str(t.Lexeme, unquoted)
}
func parseTime(p *parser, bp oper.BP, t *token.Token) ast.Expr {
	// time 字面量会被 desugar 成 strtotime, 这里留着测试场景
	ts := timelib.Strtotime(t.Lexeme[1 : len(t.Lexeme)-1]) // attach ast
	// util.Assert(ts != 0, "invalid time literal: %s", lit.Text)
	return ast.Time(t.Lexeme, ts)
}

func parseNum0(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return n, nil
	}

	if strings.HasPrefix(s, "0x") {
		n, err := strconv.ParseInt(s[2:], 16, 64)
		if err == nil {
			return float64(n), nil
		}
	}
	if strings.HasPrefix(s, "0b") {
		n, err := strconv.ParseInt(s[2:], 2, 64)
		if err == nil {
			return float64(n), nil
		}
	}
	if strings.HasPrefix(s, "0o") {
		n, err := strconv.ParseInt(s[2:], 8, 64)
		if err == nil {
			return float64(n), nil
		}
	}

	return 0, errors.New("invalid num: " + s)
}
