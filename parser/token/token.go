package token

import "github.com/goghcrow/yae/parser/loc"

type Token struct {
	Kind
	loc.Loc
	Lexeme string
}

func (t *Token) String() string { return t.Lexeme }
