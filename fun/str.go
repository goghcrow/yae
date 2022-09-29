package fun

import (
	"github.com/goghcrow/yae/token"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"unicode/utf8"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// PLUS_STR_STR + :: str -> str -> str
	// overload +
	PLUS_STR_STR = val.Fun(
		types.Fun(token.PLUS.Name(), []*types.Kind{types.Str, types.Str}, types.Str),
		func(args ...*val.Val) *val.Val {
			return val.Str(args[0].Str().V + args[1].Str().V)
		},
	)
	// LEN_STR == :: str -> num
	LEN_STR = val.Fun(
		types.Fun("len", []*types.Kind{types.Str}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(float64(utf8.RuneCountInString(args[0].Str().V)))
		},
	)
)
