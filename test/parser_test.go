package test

import (
	"fmt"
	"testing"

	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/token"
)

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
			actual := parse(tt.input, tt.ops...).String()
			if actual != tt.output {
				t.Errorf("expect %s actual %s", tt.output, actual)
			}
		})
	}
}

func TestSyntaxError(t *testing.T) {
	for _, tt := range []struct {
		s      string
		expect string
	}{
		{"a := 1", "syntax error in pos 4-1 line 1 col 4: nothing token matched"},
		{"[:}", "syntax error in pos 3-4 line 1 col 3: expect `]` actual `}`"},
		{"[1,2,3}", "syntax error in pos 2-3 line 1 col 2: expect `list or map`"},
		{`"Hello" + `, "syntax error in pos 0-0 line 0 col 0: <END-OF-FILE>"},
		{"a == b == c", "syntax error in pos 3-5 line 1 col 3: == non-infix"},              // non-infix
		{"a b", "syntax error in pos 3-4 line 1 col 3: expect `<END-OF-FILE>` actual `b`"}, // multi
	} {
		t.Run(tt.s, func(t *testing.T) {
			_, err := syntaxError(tt.s)
			if err != tt.expect {
				t.Errorf("expect syntax error `%s` actual `%s`", tt.expect, err)
			}
		})
	}
}

func syntaxError(s string) (res, err string) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Sprintf("%v", r)
		}
	}()
	res = parse(s).String()
	return
}
