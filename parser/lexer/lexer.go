package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/pos"
	"github.com/goghcrow/yae/parser/token"
)

func NewLexer(ops []oper.Operator) *lexer {
	return &lexer{
		lexicon: newLexicon(ops),
	}
}

// Lex 表达式通常都很短, 这里没有要做成语法制导按需lex, e.g. chan *token.Token
func (l *lexer) Lex(input string) []*token.Token {
	l.input = []rune(input)
	l.Pos = pos.Pos{}
	var toks []*token.Token
	for {
		t := l.next()
		if t.Kind == token.EOF {
			break
		}
		toks = append(toks, t)
	}
	return toks
}

var EOF = &token.Token{
	Kind:   token.EOF,
	Pos:    pos.Unknown,
	Lexeme: "<END-OF-FILE>",
}

type lexer struct {
	lexicon
	pos.Pos
	input []rune
}

func (l *lexer) skipSpace() {
	for l.Idx < len(l.input) {
		r := l.input[l.Idx]
		if !isSpace(r) {
			break
		}
		l.Move(r)
	}
}

func (l *lexer) next() *token.Token {
	l.skipSpace()
	if l.Idx >= len(l.input) {
		return EOF
	}

	p := l.Pos
	sub := string(l.input[l.Idx:])
	for _, rl := range l.lexicon.rules {
		offset := rl.match(sub)
		if offset >= 0 {
			matched := l.input[l.Idx : l.Idx+offset]
			for _, r := range matched {
				l.Move(r)
			}
			p.IdxEnd = l.Pos.Idx
			return &token.Token{Kind: rl.Kind, Lexeme: string(matched), Pos: p}
		}
	}
	panic(fmt.Errorf("syntax error in %s: nothing token matched", l.Pos))
}

func isSpace(r rune) bool { return unicode.IsSpace(r) }

func runeCount(s string) int { return utf8.RuneCountInString(s) }
