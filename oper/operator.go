package oper

import (
	"github.com/goghcrow/yae/token"
	"regexp"
	"strings"
)

type Operator struct {
	token.Type
	BP
	Fixity
}

const (
	// 允许自定义操作符字符列表
	operators = ":!#$%^&*+./<=>?@\\ˆ^|~-"
)

var (
	idReg = regexp.MustCompile("^[[a-zA-Z\\p{L}_][a-zA-Z0-9\\p{L}_]*]+$")
	opReg = regexp.MustCompile("^[" + regexp.QuoteMeta(operators) + "]+$")
)

func HasPrefix(s string) bool {
	for _, r := range []rune(operators) {
		if strings.HasPrefix(s, string(r)) {
			return true
		}
	}
	return false
}

func IsIdentOp(name string) bool {
	return idReg.MatchString(name)
}
