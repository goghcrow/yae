package fun

import (
	"fmt"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// PLUS_STR_STR + :: str -> str -> str
	// overload +
	PLUS_STR_STR = val.Fun(
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

	// STRING_ANY string :: forall a. (a -> str)
	STRING_ANY = func() *val.Val {
		a := types.Slot("a")
		return val.Fun(
			types.Fun(STRING, []*types.Kind{a}, types.Str),
			func(args ...*val.Val) *val.Val {
				return val.Str(stringify(args[0]))
			},
		)
	}()

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

func stringify(v *val.Val) string {
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
		return stringifyList(v.List())
	case types.TMap:
		return stringifyMap(v.Map())
	case types.TObj:
		return stringifyObj(v.Obj())
	case types.TFun:
		return "#fun"
	default:
		util.Unreachable()
		return ""
	}
}

func stringifyList(l *val.ListVal) string {
	if len(l.V) == 0 {
		return "[]"
	}

	buf := &strings.Builder{}
	buf.WriteString("[")
	isFst := true
	for _, v := range l.V {
		if isFst {
			isFst = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(stringify(v))
	}
	buf.WriteString("]")
	return buf.String()
}

func stringifyMap(m *val.MapVal) string {
	if len(m.V) == 0 {
		return "[:]"
	}

	buf := &strings.Builder{}
	buf.WriteString("[")
	isFst := true
	for k, v := range m.V {
		if isFst {
			isFst = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(k.String())
		buf.WriteString(": ")
		buf.WriteString(stringify(v))
	}
	buf.WriteString("]")
	return buf.String()
}

func stringifyObj(v *val.ObjVal) string {
	if len(v.V) == 0 {
		return "{}"
	}

	buf := &strings.Builder{}
	buf.WriteString("{")
	isFst := true
	fs := v.Kind.Obj().Fields
	for i, vl := range v.V {
		if isFst {
			isFst = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(fs[i].Name)
		buf.WriteString(": ")
		buf.WriteString(stringify(vl))
	}
	buf.WriteString("}")
	return buf.String()
}
