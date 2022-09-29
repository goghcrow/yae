package fun

import (
	"fmt"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"time"
)

//goland:noinspection GoSnakeCaseUsage
var (
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
	// TODAY_STR :: str -> time
	TODAY_STR = val.Fun(
		types.Fun("today", []*types.Kind{types.Str}, types.Time),
		func(args ...*val.Val) *val.Val {
			// 时区?
			loc := time.Local
			s := args[0].Str().V
			hm, err := time.ParseInLocation("15:04", s, loc)
			if err != nil {
				panic(fmt.Sprintf("invalid time: %s in layout hh:mm", s))
			}
			year, month, day := time.Now().Date()
			t := time.Date(year, month, day, hm.Hour(), hm.Minute(), hm.Second(), hm.Nanosecond(), loc)
			return val.Time(t)
		},
	)

	// todo yesterday() tomorrow()
)
