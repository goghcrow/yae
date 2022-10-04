package lexer

import (
	"github.com/goghcrow/yae/oper"
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
		if strings.HasPrefix(s, string(t)) {
			return len(string(t))
		} else {
			return -1
		}
	}}
}

func keyword(t token.Type) rule {
	postfix := regexp.MustCompile(`^[a-zA-Z\d\p{L}_]+`)
	return rule{t, func(s string) int {
		k := string(t) // keyword 需要匹配完整单词
		if strings.HasPrefix(s, k) && !postfix.MatchString(s[len(k):]) {
			return len(k)
		}
		return -1
	}}
}

func reg(t token.Type, pattern string) rule {
	compiled := regexp.MustCompile("^" + pattern)
	return rule{t, func(s string) int {
		found := compiled.FindString(s)
		if found == "" {
			return -1
		} else {
			return len(found)
		}
	}}
}

// primOp . ? 等内置操作符的优先级最高, 特殊处理
// e.g 比如自定义操作符 .^. 不能匹配成 [`.`, `^.`]
func primOp(t token.Type) rule {
	sz := len(string(t))
	return rule{t, func(s string) int {
		if !strings.HasPrefix(s, string(t)) {
			return -1
		}
		if len(s) == sz {
			return sz
		}
		if oper.HasPrefix(s[sz:]) {
			return -1
		} else {
			return sz
		}
	}}
}
