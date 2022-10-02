package fun

import (
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"time"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// STRTOTIME_STR strtotime :: str -> time
	STRTOTIME_STR = val.Fun(
		types.Fun("strtotime", []*types.Kind{types.Str}, types.Time),
		func(args ...*val.Val) *val.Val {
			timeStr := args[0].Str().V
			ts := util.Strtotime(timeStr)
			return val.Time(time.Unix(ts, 0))
		},
	)
)
