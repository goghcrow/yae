package timelib

import (
	"fmt"
	"os"
	"path"
	"time"
)

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
