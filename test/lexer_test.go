package test

import (
	"fmt"
	"testing"

	"github.com/goghcrow/yae/parser/lexer"
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/token"
)

func TestLoc(t *testing.T) {
	input := `if(
		布尔,
		列表[0].姓名.len() + 数字,
		0
	)`
	toks := lexer.NewLexer(oper.BuiltIn()).Lex(input)
	expect := []string{
		`if pos 1-3 line 1 col 1`,
		`( pos 3-4 line 1 col 3`,
		`布尔 pos 7-9 line 2 col 3`,
		`, pos 9-10 line 2 col 5`,
		`列表 pos 13-15 line 3 col 3`,
		`[ pos 15-16 line 3 col 5`,
		`0 pos 16-17 line 3 col 6`,
		`] pos 17-18 line 3 col 7`,
		`. pos 18-19 line 3 col 8`,
		`姓名 pos 19-21 line 3 col 9`,
		`. pos 21-22 line 3 col 11`,
		`len pos 22-25 line 3 col 12`,
		`( pos 25-26 line 3 col 15`,
		`) pos 26-27 line 3 col 16`,
		`+ pos 28-29 line 3 col 18`,
		`数字 pos 30-32 line 3 col 20`,
		`, pos 32-33 line 3 col 22`,
		`0 pos 36-37 line 4 col 3`,
		`) pos 39-40 line 5 col 2`,
	}
	for i, tok := range toks {
		actual := tok.String() + " " + tok.Pos.String()
		if expect[i] != actual {
			t.Errorf("expect %s actual %s", expect[i], actual)
		}
	}
}

func TestLexer(t *testing.T) {
	t.Run("", func(t *testing.T) {
		s := ".^."
		err := lexError(s)
		if err == nil {
			t.Errorf("expect err")
		}
	})
	t.Run("", func(t *testing.T) {
		toks := lexer.NewLexer(oper.BuiltIn()).Lex("iff")
		if len(toks) != 1 {
			t.Fail()
		}
	})
	t.Run("", func(t *testing.T) {
		toks := lexer.NewLexer([]oper.Operator{
			{
				Kind:   "as",
				BP:     oper.BP_TERM,
				Fixity: oper.INFIX_N,
			},
		}).Lex("assert")
		if len(toks) != 1 {
			t.Fail()
		}
	})
	t.Run("", func(t *testing.T) {
		s := ".^."
		toks := lexer.NewLexer(append(oper.BuiltIn(), oper.Operator{
			Kind:   token.Kind(s),
			BP:     oper.BP_TERM,
			Fixity: oper.INFIX_N,
		})).Lex(s)
		if len(toks) != 1 || toks[0].Lexeme != s {
			t.Errorf("error")
		}
	})
}

func lexError(s string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	lexer.NewLexer(oper.BuiltIn()).Lex(s)
	return nil
}
