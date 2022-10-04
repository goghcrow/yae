package lexer

import (
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
)

// lexicon Lexical grammar
type lexicon struct {
	rules []rule
}

func (l *lexicon) addRule(r ...rule) {
	l.rules = append(l.rules, r...)
}

func (l *lexicon) addOper(t token.Type) {
	if oper.IsIdentOp(string(t)) {
		l.addRule(keyword(t))
	} else {
		l.addRule(str(t))
	}
}
