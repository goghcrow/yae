package token

import (
	"strconv"
)

type Token struct {
	Type
	Lexeme string
}

func (t *Token) String() string {
	return strconv.Quote(t.Lexeme)
}
