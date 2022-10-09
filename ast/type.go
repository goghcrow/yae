package ast

type NodeType int

const (
	IDENT NodeType = iota

	LITERAL
	LIST
	MAP
	OBJ

	UNARY
	BINARY
	TENARY

	IF // if 如果是普通函数, 可以去掉
	CALL

	SUBSCRIPT
	MEMBER

	GROUP
)

func (n NodeType) String() string {
	return [...]string{
		"IDENT", "LITERAL", "LIST", "MAP", "OBJ",
		"UNARY", "BINARY", "TENARY",
		"IF",
		"CALL",
		"SUBSCRIPT", "MEMBER",
		"GROUP",
	}[n]
}

type LitType int

//goland:noinspection GoSnakeCaseUsage
const (
	LIT_NULL LitType = iota // 暂时没用
	LIT_STR
	LIT_TIME
	LIT_NUM
	LIT_TRUE
	LIT_FALSE
)

func (l LitType) String() string {
	return [...]string{"LIT_NULL", "LIT_STR", "LIT_TIME", "LIT_NUM", "LIT_TRUE", "LIT_FALSE"}[l]
}
