package fun

import (
	"github.com/goghcrow/yae/timelib"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"time"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// STRTOTIME_STR strtotime :: str -> time
	STRTOTIME_STR = val.Fun(
		types.Fun(STRTOTIME, []*types.Type{types.Str}, types.Time),
		func(args ...*val.Val) *val.Val {
			timeStr := args[0].Str().V
			ts := timelib.Strtotime(timeStr)
			return val.Time(time.Unix(ts, 0))
		},
	)
)
