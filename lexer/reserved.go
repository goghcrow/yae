package lex

// reserved 保留关键字, 以后可以改造成脚本语言
var reserved = []string{
	"byte", "int", "float", "double", "string", "bool", "boolean", "ch",
	"type", "var", "def", "define", "let", "rec", "mut", "fun", "fn", "function",
	"record", "struct", "map", "list", "object", "class", "trait", "interface",
	"sealed", "extends",
	"prefix", "infixl", "infixr", "infixn",
	"for", "do", "while", "switch", "cast", "range", "match", "select",
	"break", "continue", "return", "try", "catch", "throw", "finally",
	"import", "as", "module",
	"assert", "debugger",
}

type void struct{}

var null void = struct{}{}

var reservedSet = map[string]void{}

func init() {
	for _, name := range reserved {
		reservedSet[name] = null
	}
}

func Reserved(name string) bool {
	_, ok := reservedSet[name]
	return ok
}
