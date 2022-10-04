package oper

//goland:noinspection GoSnakeCaseUsage
const (
	PLUS  = "+"
	MINUS = "-"
	MUL   = "*"
	DIV   = "/"
	MOD   = "%"
	EXP   = "^"

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

var ops = []Operator{
	{PLUS, BP_PREFIX, PREFIX},
	{MINUS, BP_PREFIX, PREFIX},

	{PLUS, BP_TERM, INFIX_L},
	{MINUS, BP_TERM, INFIX_L},
	{MUL, BP_FACTOR, INFIX_L},
	{DIV, BP_FACTOR, INFIX_L},
	{MOD, BP_FACTOR, INFIX_L},
	{PLUS, BP_TERM, INFIX_L},
	{MINUS, BP_TERM, INFIX_L},
	{EXP, BP_EXP, INFIX_R},

	{LOGIC_OR, BP_LOGIC_OR, INFIX_L},
	{LOGIC_AND, BP_LOGIC_AND, INFIX_L},
	{LE, BP_COMP, INFIX_N},
	{LT, BP_COMP, INFIX_N},
	{GE, BP_COMP, INFIX_N},
	{GT, BP_COMP, INFIX_N},
	{EQ, BP_EQ, INFIX_N},
	{NE, BP_EQ, INFIX_N},
	{LOGIC_NOT, BP_PREFIX, PREFIX},
}

func BuildIn() []Operator {
	return ops
}
