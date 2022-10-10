package test

import (
	"fmt"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
	"testing"
)

func TestLexer(t *testing.T) {
	t.Run("", func(t *testing.T) {
		s := ".^."
		err := lexError(s)
		if err == nil {
			t.Errorf("expect err")
		}
	})
	t.Run("", func(t *testing.T) {
		toks := lexer.NewLexer(oper.BuildIn()).Lex("iff")
		if len(toks) != 1 {
			t.Fail()
		}
	})
	t.Run("", func(t *testing.T) {
		toks := lexer.NewLexer([]oper.Operator{
			{
				Type:   "as",
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
		toks := lexer.NewLexer(append(oper.BuildIn(), oper.Operator{
			Type:   token.Type(s),
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
	lexer.NewLexer(oper.BuildIn()).Lex(s)
	return nil
}
