package lexer

import (
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/token"
)

// lexicon Lexical grammar
type lexicon struct {
	rules []rule
}

func (l *lexicon) addRule(r ...rule) {
	l.rules = append(l.rules, r...)
}

func (l *lexicon) addOper(k token.Kind) {
	if oper.IsIdentOp(string(k)) {
		l.addRule(keyword(k))
	} else {
		l.addRule(str(k))
	}
}
