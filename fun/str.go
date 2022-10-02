package fun

import (
	"fmt"
	"github.com/goghcrow/yae/token"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"regexp"
	"strconv"
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
	// LEN_STR len :: str -> num
	LEN_STR = val.Fun(
		types.Fun("len", []*types.Kind{types.Str}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(float64(utf8.RuneCountInString(args[0].Str().V)))
		},
	)

	// STRING_ANY string :: forall a. (a -> str)
	STRING_ANY = func() *val.Val {
		a := types.Slot("a")
		return val.Fun(
			types.Fun("string", []*types.Kind{a}, types.Str),
			func(args ...*val.Val) *val.Val {
				return val.Str(stringfy(args[0]))
			},
		)
	}()

	// MATCH_STR_STR match :: str -> str -> bool
	MATCH_STR_STR = val.Fun(
		types.Fun("match", []*types.Kind{types.Str, types.Str}, types.Bool),
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

func stringfy(v *val.Val) string {
	switch v.Kind.Type {
	case types.TNum:
		n := v.Num()
		if n.IsInt() {
			return fmt.Sprintf("%d", n.Int())
		} else {
			return fmt.Sprintf("%f", n.V)
		}
	case types.TBool:
		return strconv.FormatBool(v.Bool().V)
	case types.TStr:
		return v.Str().V
	case types.TTime:
		return v.Time().V.String()
	case types.TList:
		return fmt.Sprintf("%s", v.List().V)
	case types.TMap:
		return fmt.Sprintf("%s", v.Map().V)
	case types.TObj:
		return fmt.Sprintf("%s", v.Obj().V)
	case types.TFun:
		return "#fun"
	default:
		util.Unreachable()
	}
	return ""
}
