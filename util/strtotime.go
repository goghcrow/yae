package util

import (
	"fmt"
	"os"
	"path"
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

const localZoneFile = "/etc/localtime" // symlinked file - set by OS

// https://stackoverflow.com/questions/68938751/how-to-get-the-full-name-of-the-local-timezone
func localTZLocationName() (name string, err error) {
	var ok bool
	if name, ok = os.LookupEnv("TZ"); ok {
		if name == "" {
			return "UTC", nil
		}
		_, err = time.LoadLocation(name) // validation
		return
	}
	fi, err := os.Lstat(localZoneFile)
	if err != nil {
		err = fmt.Errorf("failed to stat %q: %w", localZoneFile, err)
		return
	}
	if (fi.Mode() & os.ModeSymlink) == 0 {
		err = fmt.Errorf("%q is not a symlink - cannot infer name", localZoneFile)
		return
	}
	p, err := os.Readlink(localZoneFile)
	if err != nil {
		return
	}
	name, err = inferFromPath(p) // handles 1 & 2 part zone names
	return
}

//goland:noinspection SpellCheckingInspection
func inferFromPath(p string) (name string, err error) {
	dir, lname := path.Split(p)
	if len(dir) == 0 || len(lname) == 0 {
		err = fmt.Errorf("cannot infer timezone name from path: %q", p)
		return
	}
	_, fname := path.Split(dir[:len(dir)-1])
	if fname == "zoneinfo" {
		name = lname // e.g. /usr/share/zoneinfo/Japan
	} else {
		name = fname + string(os.PathSeparator) + lname // e.g. /usr/share/zoneinfo/Asia/Tokyo
	}
	_, err = time.LoadLocation(name) // validation
	return
}

// ================================================================================

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

// Strtotime
// https://www.php.net/manual/en/function.strtotime.php
// https://www.php.net/manual/en/datetime.formats.php
// https://github.com/goghcrow/strtotime
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
