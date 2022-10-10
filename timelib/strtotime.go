package timelib

import (
	"sync"
	"time"
	"unsafe"
)

/*
#include <stdlib.h>
#include "lib/timelib.h"
#include "strtotime.h"
#cgo darwin,arm64 LDFLAGS: -L./lib -ltimelib_darwin_arm64
#cgo darwin,amd64 LDFLAGS: -L./lib -ltimelib_darwin_amd64
#cgo linux,amd64 LDFLAGS: -L./lib -ltimelib_linux_amd64
*/
import "C"

var tzCache = map[string]*C.timelib_tzinfo{}
var tcCacheMut = sync.Mutex{}

//export parse_tzfile
//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
func parse_tzfile(formal_tzname *C.char, tzdb *C.timelib_tzdb, dummy_error_code *C.int) *C.timelib_tzinfo {
	g_formal_tzname := C.GoString(formal_tzname)

	tcCacheMut.Lock()
	tzi, ok := tzCache[g_formal_tzname]
	tcCacheMut.Unlock()
	if ok {
		return tzi
	}

	tzi = C.timelib_parse_tzfile(formal_tzname, tzdb, dummy_error_code)
	if tzi != nil {
		tcCacheMut.Lock()
		tzCache[g_formal_tzname] = tzi
		tcCacheMut.Unlock()
	}
	return tzi
}

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
func parseTzFile(formal_tzname string) *C.timelib_tzinfo {
	tz_name := C.CString(formal_tzname)
	defer C.free(unsafe.Pointer(tz_name))
	var dummy_error_code C.int
	return parse_tzfile(tz_name, C.timelib_builtin_db(), &dummy_error_code)
}

var localTZName = func() string {
	tzName, err := localTZLocationName()
	if err == nil {
		return tzName
	} else {
		return "UTC"
	}
}()

type StrtotimeOpt struct {
	BaseTs int64
	TzName string
}

type CallStrtotimeOpt func(opt *StrtotimeOpt)

func WithBaseTime(base time.Time) CallStrtotimeOpt {
	return func(opt *StrtotimeOpt) {
		opt.BaseTs = base.Unix()
	}
}
func WithTimeZoneString(timeZoneName string) CallStrtotimeOpt {
	return func(opt *StrtotimeOpt) {
		opt.TzName = timeZoneName
	}
}

func Strtotime(timeStr string, opts ...CallStrtotimeOpt) int64 {
	opt := &StrtotimeOpt{
		BaseTs: time.Now().Unix(),
		TzName: localTZName,
	}
	for _, set := range opts {
		set(opt)
	}
	tzi := parseTzFile(opt.TzName)
	if tzi == nil {
		return 0
	}
	times := C.CString(timeStr)
	defer C.free(unsafe.Pointer(times))
	return int64(C.strtotime(times, C.int(len(timeStr)), C.longlong(opt.BaseTs), tzi))
}
