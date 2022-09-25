package token

type Type int

//goland:noinspection GoSnakeCaseUsage
const (
	_ Type = iota

	MUL
	DIV
	MOD
	PLUS
	MINUS
	EXP
	UNARY_PLUS
	UNARY_MINUS

	GT
	LT
	LE
	GE
	EQ
	NE
	QUESTION

	LOGIC_NOT
	LOGIC_AND
	LOGIC_OR

	IF
	THEN
	ELSE
	END

	NOT
	AND
	OR

	NAME
	NUM
	STR
	NULL
	TRUE
	FALSE

	COMMA
	DOT
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACKET
	RIGHT_BRACKET
	LEFT_BRACE
	RIGHT_BRACE
	COLON

	EOF
)

const SIZE = EOF + 1

var tokens = [SIZE]struct {
	Name   string
	BP     BP
	Fixity Fixity
}{
	{"#nil", BP_NONE, NA},
	// =+=+=+=+=+=+=+=+=+=+ 算子：算术 =+=+=+=+=+=+=+=+=+=+=+
	{"*", BP_FACTOR, INFIX_L},
	{"/", BP_FACTOR, INFIX_L},
	{"%", BP_FACTOR, INFIX_L},
	{"+", BP_TERM, INFIX_L},
	{"-", BP_TERM, INFIX_L},
	{"^", BP_EXP, INFIX_R},
	{"+", BP_PREFIX_UNARY, PREFIX},
	{"-", BP_PREFIX_UNARY, PREFIX},
	// =+=+=+=+=+=+=+=+=+=+ 算子：比较 =+=+=+=+=+=+=+=+=+=+=+
	{">", BP_COMP, INFIX_L},
	{"<", BP_COMP, INFIX_L},
	{"<=", BP_COMP, INFIX_L},
	{">=", BP_COMP, INFIX_L},
	{"==", BP_EQ, INFIX_L},
	{"!=", BP_EQ, INFIX_L},
	{"?", BP_COND, INFIX_R},
	// =+=+=+=+=+=+=+=+=+=+ 算子：逻辑运算 =+=+=+=+=+=+=+=+=+=+=+
	{"!", BP_PREFIX_UNARY, INFIX_R},
	{"&&", BP_LOGIC_AND, INFIX_L},
	{"||", BP_LOGIC_OR, INFIX_L},
	// =+=+=+=+=+=+=+=+=+=+ 关键字 =+=+=+=+=+=+=+=+=+=+=+
	{"if", BP_NONE, NA},
	{"then", BP_NONE, NA},
	{"else", BP_NONE, NA},
	{"end", BP_NONE, NA},
	{"not", BP_NONE, NA},
	{"and", BP_NONE, NA},
	{"or", BP_NONE, NA},
	// =+=+=+=+=+=+=+=+=+=+ 标识符 + 字面量 =+=+=+=+=+=+=+=+=+=+=+
	{"#name", BP_NONE, NA},
	{"#num", BP_NONE, NA},
	{"#str", BP_NONE, NA},
	{"null", BP_NONE, NA},
	{"true", BP_NONE, NA},
	{"false", BP_NONE, NA},
	// =+=+=+=+=+=+=+=+=+=+ 其他 =+=+=+=+=+=+=+=+=+=+=+
	{",", BP_COMMA, INFIX_L},
	{".", BP_MEMBER, INFIX_L},
	// LED 作为CALL是左结合中缀操作符 CALL = {CALL, "(", 220, LEFT}
	// NUD 作为GROUPING是前缀操作符符号 GROUPING = {GROUPING, "(", 0, NA}
	{"(", BP_CALL, INFIX_L},
	{")", BP_NONE, NA},
	{"[", BP_MEMBER, INFIX_L},
	{"]", BP_NONE, NA},
	{"{", BP_LEFT_BRACE, INFIX_L},
	{"}", BP_NONE, NA},
	{":", BP_NONE, NA},
	{"-EOF-", BP_NONE, NA},
}

func (t Type) Bp() BP {
	return tokens[t].BP
}
func (t Type) Name() string {
	return tokens[t].Name
}
func (t Type) Fixity() Fixity {
	return tokens[t].Fixity
}
func (t Type) String() string {
	return t.Name()
}
