package fun

import (
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"regexp"
	"unicode/utf8"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// ADD_STR_STR + :: str -> str -> str
	// overload +
	ADD_STR_STR = val.Fun(
		types.Fun(oper.PLUS, []*types.Kind{types.Str, types.Str}, types.Str),
		func(args ...*val.Val) *val.Val {
			return val.Str(args[0].Str().V + args[1].Str().V)
		},
	)
	// LEN_STR len :: str -> num
	LEN_STR = val.Fun(
		types.Fun(LEN, []*types.Kind{types.Str}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(float64(utf8.RuneCountInString(args[0].Str().V)))
		},
	)

	// MATCH_STR_STR match :: str -> str -> bool
	MATCH_STR_STR = val.Fun(
		types.Fun(MATCH, []*types.Kind{types.Str, types.Str}, types.Bool),
		func(args ...*val.Val) *val.Val {
			pattern := args[0].Str().V
			s := args[1].Str().V
			matched, err := regexp.MatchString(pattern, s)
			if err != nil {
				panic(err)
			}
			return val.Bool(matched)
		},
	)
)
