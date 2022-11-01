package lexer

import "github.com/goghcrow/yae/util"

// reserved 保留关键字, 以后可以改造成脚本语言
var reserved = []string{
	"byte", "int", "float", "double", "string", "bool", "boolean", "ch", "void",
	"type", "var", "def", "define", "let", "rec", "mut", "fun", "fn", "function",
	"record", "struct", "map", "list", "object", "class", "trait", "interface",
	"sealed", "extends",
	"prefix", "infixl", "infixr", "infixn",
	"for", "do", "while", "switch", "cast", "range", "match", "select",
	"break", "continue", "return", "try", "catch", "throw", "finally",
	"import", "as", "module", "package", "namespace",
	"assert", "debugger",
}

var reservedSet = func() util.StrSet {
	s := util.StrSet{}
	for _, name := range reserved {
		s.Add(name)
	}
	return s
}()

func Reserved(name string) bool {
	return reservedSet.Contains(name)
}
