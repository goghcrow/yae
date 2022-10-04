package lexer

import (
	"fmt"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
	"testing"
)

func TestLexer(t *testing.T) {
	{
		s := ".^."
		err := lexError(s)
		if err == nil {
			t.Errorf("expect err")
		}
	}
	{
		s := ".^."
		toks := NewLexer(append(oper.BuildIn(), oper.Operator{
			Type:   token.Type(s),
			BP:     oper.BP_TERM,
			Fixity: oper.INFIX_N,
		})).Lex(s)
		if len(toks) != 1 || toks[0].Lexeme != s {
			t.Errorf("error")
		}
	}
}

func lexError(s string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	NewLexer(oper.BuildIn()).Lex(s)
	return nil
}
