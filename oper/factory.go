package oper

import "sort"

//goland:noinspection GoSnakeCaseUsage
const (
	PLUS = "+"
	SUB  = "-"
	MUL  = "*"
	DIV  = "/"
	MOD  = "%"
	EXP  = "^"

	GT = ">"
	GE = ">="
	LT = "<"
	LE = "<="
	EQ = "=="
	NE = "!="

	LOGIC_NOT = "!"
	LOGIC_AND = "&&"
	LOGIC_OR  = "||"
)

// 内置的操作符与自定义操作符一致, 可以不使用或选择使用
var ops = []Operator{
	{PLUS, BP_PREFIX, PREFIX},
	{SUB, BP_PREFIX, PREFIX}, // NEGATE

	{PLUS, BP_TERM, INFIX_L},
	{SUB, BP_TERM, INFIX_L},
	{MUL, BP_FACTOR, INFIX_L},
	{DIV, BP_FACTOR, INFIX_L},
	{MOD, BP_FACTOR, INFIX_L},
	{PLUS, BP_TERM, INFIX_L},
	{SUB, BP_TERM, INFIX_L},
	{EXP, BP_EXP, INFIX_R},

	{LE, BP_CMP, INFIX_N},
	{LT, BP_CMP, INFIX_N},
	{GE, BP_CMP, INFIX_N},
	{GT, BP_CMP, INFIX_N},
	{EQ, BP_EQ, INFIX_N},
	{NE, BP_EQ, INFIX_N},

	{LOGIC_OR, BP_LOGIC_OR, INFIX_L},
	{LOGIC_AND, BP_LOGIC_AND, INFIX_L},
	{LOGIC_NOT, BP_PREFIX, PREFIX},
}

func BuildIn() []Operator {
	return ops
}

// Sort 📢 因为 lexer 是按顺序匹配, 对于多字符的符号操作符需要注意顺序, 多字符放在单字符之前, ident 操作符不需要
// e.g. ! 需要放在 != 之后, > 需要放在 >= 之后
// e.g. 如果定义 & 需要放在  && 之后
// 使用 ops 之前, 需要先排下序
func Sort(ops []Operator) []Operator {
	sort.SliceStable(ops, func(i, j int) bool {
		x := ops[i].Type
		y := ops[j].Type
		if x == y || len(x) == len(y) {
			return false
		}
		return len(x) > len(y)
	})
	return ops
}
