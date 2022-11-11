package token

import "github.com/goghcrow/yae/parser/pos"

type Token struct {
	Kind
	pos.Pos
	Lexeme string
}

func (t *Token) String() string { return t.Lexeme }
