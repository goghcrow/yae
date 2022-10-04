package token

// Type 因为要支持动态添加操作符, int 不方便
type Type string

//goland:noinspection GoSnakeCaseUsage
const (
	QUESTION = "?"

	IF   = "if"
	THEN = "then"
	ELSE = "else"
	END  = "end"

	NAME = "'name"
	NUM  = "'num"
	STR  = "'str"
	TIME = "'time"

	NULL  = "null"
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

	EOF = "-EOF-"
)
