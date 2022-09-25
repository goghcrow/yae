package lex

import (
	"github.com/goghcrow/yae/token"
	"regexp"
	"strings"
)

// match 返回匹配 endIdx, 失败返回-1
type match func(string) int

type rule struct {
	token.Type
	match
}

func str(t token.Type) rule {
	return rule{t, func(s string) int {
		if strings.HasPrefix(s, t.Name()) {
			return len(t.Name())
		} else {
			return -1
		}
	}}
}

func keyword(t token.Type) rule {
	return keywordAlias(t, t.Name())
}

func keywordAlias(t token.Type, lexeme string) rule {
	postfix := regexp.MustCompile(`^[a-zA-Z\d\p{L}_]+`)
	return rule{t, func(s string) int {
		// keyword 需要匹配完整单词
		if strings.HasPrefix(s, lexeme) && !postfix.MatchString(s[len(lexeme):]) {
			return len(lexeme)
		}
		return -1
	}}
}

func reg(t token.Type, pattern string) rule {
	compiled := regexp.MustCompile("^" + pattern)
	return rule{t, func(sub string) int {
		found := compiled.FindString(sub)
		if found == "" {
			return -1
		} else {
			return len(found)
		}
	}}
}
