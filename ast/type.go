package ast

import "github.com/goghcrow/yae/util"

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
	switch l {
	case LIT_NULL:
		return "lit_num"
	case LIT_STR:
		return "lit_str"
	case LIT_TIME:
		return "lit_time"
	case LIT_NUM:
		return "lit_num"
	case LIT_TRUE:
		return "lit_true"
	case LIT_FALSE:
		return "lit_false"
	default:
		util.Unreachable()
		return ""
	}
}
