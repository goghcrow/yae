package lexer

import (
	"regexp"
	"strings"

	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/token"
)

const NotMatched = -1

type rule struct {
	token.Type
	match func(string) int // 匹配返回 EndRuneCount , 失败返回 NotMatched
}

func str(t token.Type) rule {
	tok := string(t)
	return rule{t, func(s string) int {
		if strings.HasPrefix(s, tok) {
			return runeCount(tok)
		} else {
			return NotMatched
		}
	}}
}

var keywordPostfix = regexp.MustCompile(`^[a-zA-Z\d\p{L}_]+`)

func keyword(t token.Type) rule {
	kw := string(t)
	return rule{t, func(s string) int {
		completedWord := strings.HasPrefix(s, kw) &&
			!keywordPostfix.MatchString(s[len(kw):])
		if completedWord {
			return runeCount(kw)
		} else {
			return NotMatched
		}
	}}
}

func reg(t token.Type, pattern string) rule {
	startWith := regexp.MustCompile("^" + pattern)
	return rule{t, func(s string) int {
		found := startWith.FindString(s)
		if found == "" {
			return NotMatched
		} else {
			return runeCount(found)
		}
	}}
}

// primOper . ? 内置操作符的优先级高于自定义操作符, 且不是匹配最长, 需要特殊处理
// e.g 比如自定义操作符 .^. 不能匹配成 [`.`, `^.`]
func primOper(t token.Type) rule {
	op := string(t)
	return rule{t, func(s string) int {
		if !strings.HasPrefix(s, op) {
			return NotMatched
		}
		completedOper := len(s) == len(op) || !oper.HasPrefix(s[len(op):])
		if completedOper {
			return runeCount(op)
		} else {
			return NotMatched
		}
	}}
}
