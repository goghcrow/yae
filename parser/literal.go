package parser

import (
	"errors"
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/timelib"
	"github.com/goghcrow/yae/util"
	"strconv"
	"strings"
)

func parseLit(typ ast.LitType, lexeme string) *ast.Expr {
	expr := ast.Literal(typ, lexeme)
	lit := expr.Literal()
	var err error

	switch typ {
	case ast.LIT_STR:
		lit.Val, err = strconv.Unquote(lit.Text) // attach ast
		util.Assert(err == nil, "invalid string literal: %s", lit.Text)
	case ast.LIT_TIME:
		// time 字面量会被 desugar 成 strtotime, 这里留着测试场景
		lit.Val = timelib.Strtotime(lit.Text[1 : len(lit.Text)-1]) // attach ast
		// util.Assert(ts != 0, "invalid time literal: %s", lit.Text)
	case ast.LIT_NUM:
		lit.Val, err = parseNum(lit.Text) // attach ast
		util.Assert(err == nil, "invalid num literal %s", lit.Text)
	case ast.LIT_TRUE:
		lit.Val = true // attach ast
	case ast.LIT_FALSE:
		lit.Val = false // attach ast
	default:
		util.Unreachable()
		return nil
	}

	return expr
}

func parseNum(s string) (float64, error) {
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
