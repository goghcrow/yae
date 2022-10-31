package token

import "github.com/goghcrow/yae/loc"

type Token struct {
	Type
	loc.Loc
	Lexeme string
}

func (t *Token) String() string {
	// return strconv.Quote(t.Lexeme)
	return t.Lexeme
}
