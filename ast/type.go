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

	IF
	CALL

	SUBSCRIPT
	MEMBER

	BEGIN
)

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
