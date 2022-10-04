package lexer

import (
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
)

func NewLexer(ops []oper.Operator) *lexer {
	return &lexer{
		lexicon: newLexicon(ops),
	}
}

func (l *lexer) Lex(input string) []*token.Token {
	l.input = input
	l.idx = 0
	toks := make([]*token.Token, 0)
	for {
		t := l.next()
		if t == EOF {
			break
		}
		toks = append(toks, t)
	}
	return toks
}

var EOF = &token.Token{Type: token.EOF}

type lexer struct {
	lexicon
	input string
	idx   int
}

func (l *lexer) skipSpace() {
	isSpace := func(c byte) bool {
		return c == ' ' || c == '\t' || c == '\r' || c == '\n'
	}
	for l.idx < len(l.input) {
		if !isSpace(l.input[l.idx]) {
			break
		}
		l.idx++
	}
}

func (l *lexer) next() *token.Token {
	l.skipSpace()

	if l.idx >= len(l.input) {
		return EOF
	}

	sub := l.input[l.idx:]
	for _, r := range l.lexicon.rules {
		offset := r.match(sub)
		if offset >= 0 {
			matched := l.input[l.idx : l.idx+offset]
			l.idx += offset
			return &token.Token{Type: r.Type, Lexeme: matched}
		}
	}
	panic("syntax error: " + sub)
}
