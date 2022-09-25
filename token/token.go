package token

import "fmt"

type Token struct {
	Type   Type
	Lexeme string
}

func (t *Token) String() string {
	return fmt.Sprintf("%q", t.Lexeme)
}
