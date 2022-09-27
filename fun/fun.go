package fun

import (
	"fmt"
	"github.com/goghcrow/yae/token"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"math"
	"time"
)

// build-in funs

// 📢 命名规则, 函数或者符号名称[_参数1类型_参数2类型..._参数n类型]

// 关于 if :: forall α. (bool -> α -> a -> α)
// if 也可以直接定义成普通 Lazy 函数, 而不需要在 typecheck 和 compile 特殊处理
// but, 因为没做类型变量, 没做 ∀, 所以需要重载成 if[bool,num,num], if[bool,str,str]... 复制很多份
// 或者在 hack 一个 top 类型, 定义成 IF_BOOL_ANY_ANY  ::  bool -> any -> any -> any
// 编译期方法分派时候时候, 实际类型没找到回退到 any 类型签名查找, 实际上就会引入子类型了, 还错误的把参数协变了
// so, 在 typecheck 和 compile 把 if 特殊处理了, and or 之类做成普通参数惰性求值函数

// 考虑下 数学计算可以先转换成 big.Int/big.Float 进行计算在转回来
// 表达式的场景对精度要求可能也那么高

//goland:noinspection GoUnusedGlobalVariable,GoSnakeCaseUsage
var (
	AnyObj = types.Obj(map[string]*types.Kind{})

	// IF_BOOL_A_A + :: forall a. (bool -> α -> α -> α)
	// if 可以声明成惰性求值的泛型函数
	IF_BOOL_A_A = func() *val.Val {
		T := types.Slot("a")
		return val.LazyFun(
			types.Fun(token.IF.Name(), []*types.Kind{types.Bool, T, T}, T),
			// 注意 if 是 lazyFun, 参数都是 thunk
			func(args ...*val.Val) *val.Val {
				if args[0].Fun().V().Bool().V {
					return args[1].Fun().V()
				} else {
					return args[2].Fun().V()
				}
			},
		)
	}()

	// PLUS_NUM + :: num -> num
	PLUS_NUM = val.Fun(
		types.Fun(token.PLUS.Name(), []*types.Kind{types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return args[0]
		},
	)
	// PLUS_NUM_NUM + :: num -> num -> num
	PLUS_NUM_NUM = val.Fun(
		types.Fun(token.PLUS.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V + args[1].Num().V)
		},
	)
	// PLUS_STR_STR + :: str -> str -> str
	PLUS_STR_STR = val.Fun(
		types.Fun(token.PLUS.Name(), []*types.Kind{types.Str, types.Str}, types.Str),
		func(args ...*val.Val) *val.Val {
			return val.Str(args[0].Str().V + args[1].Str().V)
		},
	)

	// MINUS_NUM - :: num -> num
	MINUS_NUM = val.Fun(
		types.Fun(token.MINUS.Name(), []*types.Kind{types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(-args[0].Num().V)
		},
	)
	// MINUS_NUM_NUM - :: num -> num -> num
	MINUS_NUM_NUM = val.Fun(
		types.Fun(token.MINUS.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V - args[1].Num().V)
		},
	)
	// MINUS_TIME_TIME - :: time -> time -> num
	MINUS_TIME_TIME = val.Fun(
		types.Fun(token.MINUS.Name(), []*types.Kind{types.Time, types.Time}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Time().V.Sub(args[1].Time().V).Seconds())
		},
	)

	// MUL_NUM_NUM * :: num -> num -> num
	MUL_NUM_NUM = val.Fun(
		types.Fun(token.MUL.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V * args[1].Num().V)
		},
	)
	// DIV_NUM_NUM / :: num -> num -> num
	DIV_NUM_NUM = val.Fun(
		types.Fun(token.DIV.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V / args[1].Num().V)
		},
	)
	// MOD_NUM_NUM % :: num -> num -> num
	MOD_NUM_NUM = val.Fun(
		types.Fun(token.MOD.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(float64(int64(args[0].Num().V) % int64(args[1].Num().V)))
		},
	)
	// EXP_NUM_NUM ^ :: num -> num -> num
	EXP_NUM_NUM = val.Fun(
		types.Fun(token.EXP.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Pow(args[0].Num().V, args[1].Num().V))
		},
	)

	//AND_BOOL_BOOL - :: bool -> bool -> bool
	AND_BOOL_BOOL = val.LazyFun(
		types.Fun(token.LOGIC_AND.Name(), []*types.Kind{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			thunk1 := args[0].Fun()
			if thunk1.V().Bool().V {
				thunk2 := args[1].Fun()
				return thunk2.V().Bool().Vl()
			} else {
				return val.False
			}
		},
	)
	//OR_BOOL_BOOL - :: bool -> bool -> bool
	OR_BOOL_BOOL = val.LazyFun(
		types.Fun(token.LOGIC_OR.Name(), []*types.Kind{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			thunk1 := args[0].Fun()
			if thunk1.V().Bool().V {
				return val.True
			} else {
				thunk2 := args[1].Fun()
				return thunk2.V().Bool().Vl()
			}
		},
	)
	// NOT_BOOL - :: bool -> bool
	NOT_BOOL = val.Fun(
		types.Fun(token.LOGIC_OR.Name(), []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(!args[0].Bool().V)
		},
	)

	// EQ_BOOL_BOOL == :: bool -> bool -> bool
	EQ_BOOL_BOOL = val.Fun(
		types.Fun(token.EQ.Name(), []*types.Kind{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Bool().V == args[1].Bool().V)
		},
	)
	// NE_BOOL_BOOL == :: bool -> bool -> bool
	NE_BOOL_BOOL = val.Fun(
		types.Fun(token.NE.Name(), []*types.Kind{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Bool().V != args[1].Bool().V)
		},
	)
	// EQ_NUM_NUM == :: num -> num -> bool
	EQ_NUM_NUM = val.Fun(
		types.Fun(token.EQ.Name(), []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(val.Equals(args[0], args[1]))
		},
	)
	// NE_NUM_NUM != :: num -> num -> bool
	NE_NUM_NUM = val.Fun(
		types.Fun(token.NE.Name(), []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(!val.Equals(args[0], args[1]))
		},
	)
	// EQ_STR_STR == :: str -> str -> bool
	EQ_STR_STR = val.Fun(
		types.Fun(token.EQ.Name(), []*types.Kind{types.Str, types.Str}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Str().V != args[1].Str().V)
		},
	)
	// NE_STR_STR == :: str -> str -> bool
	NE_STR_STR = val.Fun(
		types.Fun(token.NE.Name(), []*types.Kind{types.Str, types.Str}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Str().V != args[1].Str().V)
		},
	)
	// EQ_TIME_TIME == :: time -> time -> bool
	EQ_TIME_TIME = val.Fun(
		types.Fun(token.EQ.Name(), []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.Equal(args[1].Time().V))
		},
	)
	// NE_TIME_TIME == :: time -> time -> bool
	NE_TIME_TIME = val.Fun(
		types.Fun(token.NE.Name(), []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(!args[0].Time().V.Equal(args[1].Time().V))
		},
	)
	//EQ_LIST_LIST == :: forall a. (list[a] -> list[a] -> bool)
	EQ_LIST_LIST = func() *val.Val {
		T := types.Slot("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(token.EQ.Name(), []*types.Kind{listT, listT}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(val.Equals(args[0], args[1]))
			},
		)
	}()
	//NE_LIST_LIST != :: forall a. (list[a] -> list[a] -> bool)
	NE_LIST_LIST = func() *val.Val {
		T := types.Slot("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(token.NE.Name(), []*types.Kind{listT, listT}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(!val.Equals(args[0], args[1]))
			},
		)
	}()
	//EQ_MAP_MAP == :: forall k v . (map[k,v] -> map[k,v] -> bool)
	EQ_MAP_MAP = func() *val.Val {
		K := types.Slot("k")
		V := types.Slot("v")
		mapKV := types.Map(K, V)
		return val.Fun(
			types.Fun(token.EQ.Name(), []*types.Kind{mapKV, mapKV}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(val.Equals(args[0], args[1]))
			},
		)
	}()
	//NE_MAP_MAP != :: forall k v . (map[k,v] -> map[k,v] -> bool)
	NE_MAP_MAP = func() *val.Val {
		K := types.Slot("k")
		V := types.Slot("v")
		mapKV := types.Map(K, V)
		return val.Fun(
			types.Fun(token.NE.Name(), []*types.Kind{mapKV, mapKV}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(!val.Equals(args[0], args[1]))
			},
		)
	}()

	// GT_NUM_NUM > :: num -> num -> bool
	GT_NUM_NUM = val.Fun(
		types.Fun(token.GT.Name(), []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Num().V > args[1].Num().V && !val.Equals(args[0], args[1]))
		},
	)
	// GE_NUM_NUM >= :: num -> num -> bool
	GE_NUM_NUM = val.Fun(
		types.Fun(token.GE.Name(), []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Num().V > args[1].Num().V || val.Equals(args[0], args[1]))
		},
	)
	// LT_NUM_NUM < :: num -> num -> bool
	LT_NUM_NUM = val.Fun(
		types.Fun(token.LT.Name(), []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Num().V < args[1].Num().V && !val.Equals(args[0], args[1]))
		},
	)
	// LE_NUM_NUM <= :: num -> num -> bool
	LE_NUM_NUM = val.Fun(
		types.Fun(token.LE.Name(), []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Num().V < args[1].Num().V || val.Equals(args[0], args[1]))
		},
	)

	// GT_TIME_TIME > :: time -> time -> bool
	GT_TIME_TIME = val.Fun(
		types.Fun(token.GT.Name(), []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.After(args[1].Time().V))
		},
	)
	// GE_TIME_TIME >= :: time -> time -> bool
	GE_TIME_TIME = val.Fun(
		types.Fun(token.GE.Name(), []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.After(args[1].Time().V) || args[0].Time().V.Equal(args[1].Time().V))
		},
	)
	// LT_TIME_TIME < :: time -> time -> bool
	LT_TIME_TIME = val.Fun(
		types.Fun(token.LT.Name(), []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.Before(args[1].Time().V))
		},
	)
	// LE_TIME_TIME <= :: time -> time -> bool
	LE_TIME_TIME = val.Fun(
		types.Fun(token.LE.Name(), []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.Before(args[1].Time().V) || args[0].Time().V.Equal(args[1].Time().V))
		},
	)

	// TIME_STR_STR time :: str -> str -> time
	TIME_STR_STR = val.Fun(
		types.Fun("time", []*types.Kind{types.Str, types.Str}, types.Time),
		func(args ...*val.Val) *val.Val {
			// 时区?
			layout := args[0].Str().V
			s := args[1].Str().V
			t, err := time.ParseInLocation(layout, s, time.Local)
			if err != nil {
				panic(fmt.Sprintf("invalid time: %s in layout %s", s, layout))
			}
			return val.Time(t)
		},
	)
	// NOW :: time
	NOW = val.Fun(
		types.Fun("now", []*types.Kind{}, types.Time),
		func(args ...*val.Val) *val.Val {
			// 时区? 可以重载一个带时区参数的函数
			return val.Time(time.Now())
		},
	)
	// TODAY :: time
	TODAY = val.Fun(
		types.Fun("today", []*types.Kind{}, types.Time),
		func(args ...*val.Val) *val.Val {
			// 时区?
			year, month, day := time.Now().Date()
			today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
			return val.Time(today)
		},
	)
	// TODAY_NUM :: num -> time
	TODAY_NUM = val.Fun(
		types.Fun("today", []*types.Kind{types.Num}, types.Time),
		func(args ...*val.Val) *val.Val {
			// 时区?
			hm := int64(args[0].Num().V)
			hh := hm / 100
			mm := hm % 100
			year, month, day := time.Now().Date()
			today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
			return val.Time(today.Add(time.Duration(hh)*time.Hour + time.Duration(mm)*time.Minute))
		},
	)

	// 目前重载比较简单, 不定参数+重载的分派会变复杂, 暂时可以用重载 n 个参数来解决, 会比较啰嗦, e.g. MAX_NUM_NUM_NUM

	// MAX_NUM_NUM + :: num -> num -> num
	MAX_NUM_NUM = val.Fun(
		types.Fun("max", []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Max(args[0].Num().V, args[1].Num().V))
		},
	)
	// MAX_NUM_NUM_NUM + :: num -> num -> num -> num
	MAX_NUM_NUM_NUM = val.Fun(
		types.Fun("max", []*types.Kind{types.Num, types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(
				math.Max(
					math.Max(args[0].Num().V, args[1].Num().V),
					args[2].Num().V,
				),
			)
		},
	)
	// MIN_NUM_NUM + :: num -> num -> num
	MIN_NUM_NUM = val.Fun(
		types.Fun("min", []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Min(args[0].Num().V, args[1].Num().V))
		},
	)
	// MIN_NUM_NUM_NUM + :: num -> num -> num -> num
	MIN_NUM_NUM_NUM = val.Fun(
		types.Fun("min", []*types.Kind{types.Num, types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(
				math.Min(
					math.Min(args[0].Num().V, args[1].Num().V),
					args[2].Num().V,
				),
			)
		},
	)
)
