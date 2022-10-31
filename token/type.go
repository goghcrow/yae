package token

// Type 因为要支持动态添加操作符, 所以 Type 没有定义成 int 枚举
// 这里需要自己保证 type 值不重复
type Type string

// '前缀的为 psuido
//
//goland:noinspection GoSnakeCaseUsage
const (
	QUESTION = "?"

	// if 是普通函数, 这里不需要
	//IF = "if"
	//THEN = "then"
	//ELSE = "else"
	//END  = "end"

	NAME = "'name"
	NUM  = "'num"
	STR  = "'str"
	TIME = "'time"

	TRUE  = "true"
	FALSE = "false"

	COMMA = ","
	DOT   = "."

	LEFT_PAREN    = "("
	RIGHT_PAREN   = ")"
	LEFT_BRACKET  = "["
	RIGHT_BRACKET = "]"
	LEFT_BRACE    = "{"
	RIGHT_BRACE   = "}"

	COLON = ":"

	EOF = "'EOF"
)
