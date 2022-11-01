package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/goghcrow/yae/parser/loc"
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/token"
)

func NewLexer(ops []oper.Operator) *lexer {
	return &lexer{
		lexicon: newLexicon(oper.Sort(ops)),
	}
}

// Lex 表达式通常都很短, 这里没有要做成语法制导按需lex, e.g. chan *token.Token
func (l *lexer) Lex(input string) []*token.Token {
	l.input = []rune(input)
	l.Loc = loc.Loc{}
	var toks []*token.Token
	for {
		t := l.next()
		if t.Type == token.EOF {
			break
		}
		toks = append(toks, t)
	}
	return toks
}

var EOF = &token.Token{
	Type:   token.EOF,
	Loc:    loc.Unknown,
	Lexeme: "<END-OF-FILE>",
}

type lexer struct {
	lexicon
	loc.Loc
	input []rune
}

func (l *lexer) skipSpace() {
	for l.Pos < len(l.input) {
		r := l.input[l.Pos]
		if !isSpace(r) {
			break
		}
		l.Move(r)
	}
}

func (l *lexer) next() *token.Token {
	l.skipSpace()
	if l.Pos >= len(l.input) {
		return EOF
	}

	pos := l.Loc
	sub := string(l.input[l.Pos:])
	for _, rl := range l.lexicon.rules {
		offset := rl.match(sub)
		if offset >= 0 {
			matched := l.input[l.Pos : l.Pos+offset]
			for _, r := range matched {
				l.Move(r)
			}
			pos.PosEnd = l.Loc.Pos
			return &token.Token{Type: rl.Type, Lexeme: string(matched), Loc: pos}
		}
	}
	panic(fmt.Errorf("syntax error in %s: nothing token matched", l.Loc))
}

func isSpace(r rune) bool { return unicode.IsSpace(r) }

func runeCount(s string) int { return utf8.RuneCountInString(s) }
