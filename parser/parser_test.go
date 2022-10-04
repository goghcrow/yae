package parser

import (
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/trans"
	"testing"
)

func parse(s string, ops ...oper.Operator) string {
	ops = append(oper.BuildIn(), ops...)
	toks := lexer.NewLexer(ops).Lex(s)
	ast := NewParser(ops).Parse(toks)
	return trans.Desugar(ast).String()
}

func TestParser(t *testing.T) {
	tests := []struct {
		input  string
		output string
		ops    []oper.Operator
	}{
		// call
		{"a()", "a()", nil},
		{"a(b)", "a(b)", nil},
		{"a(b, c)", "a(b, c)", nil},
		{"a(b)(c)", "a(b)(c)", nil},
		{"a(b) + c(d)", "+(a(b), c(d))", nil},
		{"a(b ? c : d, e + f)", "a(if(b, c, d), +(e, f))", nil},

		// Unary precedence
		{"-1", "-(1)", nil},
		{"+1", "+(1)", nil},
		{"~!-+a", "~(!(-(+(a))))", []oper.Operator{{
			Type:   token.Type("~"),
			BP:     oper.BP_PREFIX,
			Fixity: oper.PREFIX,
		}}},
		{"a!~!", "!(~(!(a)))", []oper.Operator{{
			Type:   token.Type("!"),
			BP:     oper.BP_POSTFIX,
			Fixity: oper.POSTFIX,
		}, {
			Type:   token.Type("~"),
			BP:     oper.BP_POSTFIX,
			Fixity: oper.POSTFIX,
		}}},

		// Unary and binary precedence
		{"-42 == 1", "==(-(42), 1)", nil},
		{"1 + -1", "+(1, -(1))", nil},
		{"1 - -1", "-(1, -(1))", nil},
		{"-a * b", "*(-(a), b)", nil},
		{"!a + b", "+(!(a), b)", nil},
		{"~a ^ b", "^(~(a), b)", []oper.Operator{{
			Type:   token.Type("~"),
			BP:     oper.BP_PREFIX,
			Fixity: oper.PREFIX,
		}}},
		{"-a!", "-(!(a))", []oper.Operator{{
			Type:   token.Type("!"),
			BP:     oper.BP_POSTFIX,
			Fixity: oper.POSTFIX,
		}}},
		{"!a#", "!(#(a))", []oper.Operator{{
			Type:   token.Type("#"),
			BP:     oper.BP_POSTFIX,
			Fixity: oper.POSTFIX,
		}}},
		{"-1 # 2", "#(-(1), 2)", []oper.Operator{{
			Type:   token.Type("#"),
			BP:     oper.BP_PREFIX,
			Fixity: oper.INFIX_L,
		}}},
		{"-1 # 2", "-(#(1, 2))", []oper.Operator{{
			Type:   token.Type("#"),
			BP:     oper.BP_PREFIX + 1,
			Fixity: oper.INFIX_L,
		}}},

		// Binary precedence
		{"a == b + c * d ^ e - f / g", "==(a, -(+(b, *(c, ^(d, e))), /(f, g)))", nil},
		{"1 - 2 + 3 * 4", "+(-(1, 2), *(3, 4))", nil},

		// Binary associativity
		{"a + b - c", "-(+(a, b), c)", nil},
		{"a * b / c", "/(*(a, b), c)", nil},
		{"a ^ b ^ c", "^(a, ^(b, c))", nil},

		// Conditional operator
		{"a ? b : c ? d : e", "if(a, b, if(c, d, e))", nil},
		{"a ? b ? c : d : e", "if(a, if(b, c, d), e)", nil},
		{"a + b ? c * d : e / f", "if(+(a, b), *(c, d), /(e, f))", nil},

		// Grouping
		{"a + (b + c) + d", "+(+(a, +(b, c)), d)", nil},
		{"a ^ (b + c)", "^(a, +(b, c))", nil},
		{"(!a)@", "@(!(a))", []oper.Operator{{
			Type:   token.Type("@"),
			BP:     oper.BP_POSTFIX,
			Fixity: oper.POSTFIX,
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := parse(tt.input, tt.ops...)
			if actual != tt.output {
				t.Errorf("expect %s actual %s", tt.output, actual)
			}
		})
	}
}

func TestSyntaxError(t *testing.T) {
	for _, expr := range []string{
		`"Hello" + `,
		"a == b == c", // non-infix
		"a b",         // multi
	} {
		s := syntaxError(expr)
		if s != "" {
			t.Errorf("expect syntax error actual `%s`", expr)
		}
	}
}

func syntaxError(s string) (r string) {
	defer func() {
		if r := recover(); r != nil {
			r = ""
		}
	}()
	return parse(s)
}
