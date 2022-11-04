package token

// Kind 因为要支持动态添加操作符, 所以 Kind 没有定义成 int 枚举
// 这里需要自己保证 Kind 值不重复
type Kind string

//goland:noinspection GoSnakeCaseUsage
const (
	TRUE          = "true"
	FALSE         = "false"
	QUESTION      = "?"
	COMMA         = ","
	DOT           = "."
	COLON         = ":"
	LEFT_PAREN    = "("
	RIGHT_PAREN   = ")"
	LEFT_BRACKET  = "["
	RIGHT_BRACKET = "]"
	LEFT_BRACE    = "{"
	RIGHT_BRACE   = "}"
	SYM           = "<sym>"
	NUM           = "<num>"
	STR           = "<str>"
	TIME          = "<time>"
	EOF           = "<END-OF-FILE>"

	// if 是普通函数, 这里不需要
	//IF = "if"
	//THEN = "then"
	//ELSE = "else"
	//END  = "end"
)
