# timelib

PHP & MongoDB 的时间解析库的 CGO wrapper. 主要为了移植 `strtotime` 函数.

## strtotime 文档

> strtotime — Parse about any English textual datetime description into a Unix timestamp

[strtotime](https://github.com/goghcrow/strtotime)

[datetime.formats](https://www.php.net/manual/en/datetime.formats.php)

[之前移植的 java 版本](https://github.com/goghcrow/strtotime)

## strtotime 示例

```golang
package timelib

import (
	"testing"
	"time"
)

func Test_Strtotime_Example(t *testing.T) {
	t.Run("母亲节", func(t *testing.T) {
		t.Log(time.Unix(Strtotime("second sunday of may 2022"), 0))
	})

	t.Run("父亲节", func(t *testing.T) {
		t.Log(time.Unix(Strtotime("third sunday of june 2022"), 0))
	})

	t.Run("感恩节", func(t *testing.T) {
		t.Log(time.Unix(Strtotime("fourth thursday of november 2022"), 0))
	})

	t.Run("2022年二月的最后一天 12:22", func(t *testing.T) {
		t.Log(time.Unix(Strtotime("last day of february 2022 12:22"), 0))
	})
	
	times := []string{
		"now",
		"yesterday",
		"today",
		"tomorrow",
		"noon",
		"midnight",

		"yesterday 08:15pm",
		"yesterday noon",
		"yesterday midnight",
		"tomorrow 18:00",
		"tomorrow moon",

		"+1 week 2 days 4 hours 2 seconds",

		"saturday this week",

		"next year",
		"next month",

		"last day",
		"last wed",

		"this week",
		"next week",
		"last week",
		"previous week",

		"monday",
		"mon",
		"tuesday",
		"tue",
		"wednesday",
		"wed",
		"thursday",
		"thu",
		"friday",
		"fri",
		"saturday",
		"sat",
		"sunday",
		"sun",

		"first day",
		"first day next month",
		"first day of next month",
		"last day next month",
		"last day of next month",
		"last day of april",

		"third Monday December 2020",
		"second Friday Nov 2022",
		"+3 week Thursday Nov 2020",
		"last wednesday of march 2020",

		"2020W30",

		"2020W101T05:00+0",

		"10/22/1990",
		"10/22",
		"01/01",

		"Sun 2020-01-01",
		"Mon 2020-01-02",

		"19970523091528",
		"20001231185859",
		"20800410101010",

		"Fri 2020-01-06",

		"2020-06-25 14:18:48.543728 America/New_York",

		"2020-10-22 13:00:00 Asia/Shanghai",
		"2022-01-01 13:00:00 UTC",
		"2020-01-01 00:00:00 Europe/Rome",

		"2020-11-26T18:51:44+01:00",
		"Thursday, 26-Nov-2020 18:51:44 CET",
		"2020-11-26T18:51:44+0100",
		"Thu, 26 Nov 20 18:51:44 +0100",
		"Thursday, 26-Nov-20 18:51:44 CET",
		"Thu, 26 Nov 2020 18:51:44 +0100",

		"May 18th 2020 5:05pm",
		"2005-8-12",
		"Sat 26th Nov 2020 18:18",
		"26th Nov",
		"Dec. 4th, 2020",
		"December 4th, 2020",
		"Sun, 13 Nov 2020 22:56:10 -0800 (PST)",
		"May 18th 5:05pm",
	}

	for _, it := range times {
		t.Log(time.Unix(Strtotime(it), 0))
	}
}
```

## 支持解析的格式

📢📢📢 **支持以下所有格式组合使用，注意相同元素不能重复设置，比如时区不能设置两次**

```java
reAgo       = "^ago";

reHour24    = "(2[0-4]|[01]?[0-9])";
reHour24Lz  = "([01][0-9]|2[0-4])";
reHour12    = "(1[0-2]|0?[1-9])";
reMinute    = "([0-5]?[0-9])";
reMinuteLz  = "([0-5][0-9])";
reSecond    = "([0-5]?[0-9]|60)";
reSecondLz  = "([0-5][0-9]|60)";
reFrac      = "(?:\\.([0-9]+))";

reMeridian      = "(?:([ap])\\.?m\\.?(?:[ \\t]|$))";

reYear          = "([0-9]{1,4})";
reYear2         = "([0-9]{2})";
reYear4         = "([0-9]{4})";
reYear4WithSign = "([+-]?[0-9]{4})";

reMonth         = "(1[0-2]|0?[0-9])";
reMonthLz       = "(0[0-9]|1[0-2])";

reMonthFull     = "january|february|march|april|may|june|july|august|september|october|november|december";
reMonthAbbr     = "jan|feb|mar|apr|may|jun|jul|aug|sept?|oct|nov|dec";
reMonthRoman    = "i{1,3}|i[vx]|vi{0,3}|xi{0,2}";
reMonthText     = '(' + reMonthFull + '|' + reMonthAbbr + '|' + reMonthRoman + ')';

reDay   = "(?:([0-2]?[0-9]|3[01])(?:st|nd|rd|th)?)";
reDayLz = "(0[0-9]|[1-2][0-9]|3[01])";

reDayFull = "sunday|monday|tuesday|wednesday|thursday|friday|saturday";
reDayAbbr = "sun|mon|tue|wed|thu|fri|sat";
reDayText = reDayFull + '|' + reDayAbbr + '|' + "weekdays?";

reDayOfYear = "(00[1-9]|0[1-9][0-9]|[12][0-9][0-9]|3[0-5][0-9]|36[0-6])";
reWeekOfYear = "(0[1-9]|[1-4][0-9]|5[0-3])";

reTzCorrection = "((?:GMT)?([+-])" + reHour24 + ":?" + reMinute + "?)";
reTzAbbr = "\\(?([a-zA-Z]{1,6})\\)?";
reTz = "[A-Z][a-z]+([_/-][A-Za-z_]+)+|" + reTzAbbr;


/* Time formats */
reTimeTiny12  = '^' + reHour12                                           + reSpaceOpt + reMeridian;
reTimeShort12 = '^' + reHour12 + "[:.]" + reMinuteLz + reSpaceOpt + reMeridian;
reTimeLong12  = '^' + reHour12 + "[:.]" + reMinute + "[:.]" + reSecondLz + reSpaceOpt + reMeridian;

reTimeShort24 = "^t?" + reHour24 + "[:.]" + reMinute;
reTimeLong24  = "^t?" + reHour24 + "[:.]" + reMinute + "[:.]" + reSecond;
reISO8601Long = "^t?" + reHour24 + "[:.]" + reMinute + "[:.]" + reSecond + reFrac;

reTzText = '(' + reTzCorrection + '|' + reTz + ')';

reISO8601NormTz = "^t?" + reHour24 + "[:.]" + reMinute + "[:.]" + reSecondLz + reSpaceOpt + reTzText;

/* gnu */
reGNUNoColon = "^t?" + reHour24Lz + reMinuteLz;
reISO8601NoColon = "^t?" + reHour24Lz + reMinuteLz + reSecondLz;

/* Date formats */
reAmericanShort     = '^' + reMonth + '/' + reDay;
reAmerican          = '^' + reMonth + '/' + reDay + '/' + reYear;
reISO8601DateSlash  = '^' + reYear4 + '/' + reMonthLz + '/' + reDayLz + "/?";
reDateSlash         = '^' + reYear4 + '/' + reMonth + '/' + reDay;
reISO8601Date4      = '^' + reYear4WithSign + '-' + reMonthLz + '-' + reDayLz;
reISO8601Date2      = '^' + reYear2 + '-' + reMonthLz + '-' + reDayLz;
reGNUDateShorter    = '^' + reYear4 + '-' + reMonth;
reGNUDateShort      = '^' + reYear + '-' + reMonth + '-' + reDay;
rePointedDate4      = '^' + reDay + "[.\\t-]" + reMonth + "[.-]" + reYear4;
rePointedDate2      = '^' + reDay + "[.\\t]" +  reMonth + "\\." + reYear2;
reDateFull          = '^' + reDay + "[ \\t.-]*" + reMonthText + "[ \\t.-]*" + reYear;
reDateNoDay         = '^' + reMonthText + "[ .\\t-]*" + reYear4;
reDateNoDayRev      = '^' + reYear4 + "[ .\\t-]*" + reMonthText;
reDateTextual       = '^' + reMonthText + "[ .\\t-]*" + reDay + "[,.stndrh\\t ]+" + reYear;
reDateNoYear        = '^' + reMonthText + "[ .\\t-]*" + reDay + "[,.stndrh\\t ]*";
reDateNoYearRev     = '^' + reDay + "[ .\\t-]*" + reMonthText;
reDateNoColon       = '^' + reYear4 + reMonthLz + reDayLz;

/* Special formats */
// 参见 https://www.php.net/manual/en/datetime.formats.compound.php
// 木有遵守这个：The "T" in the SOAP, XMRPC and WDDX formats is case-sensitive, you can only use the upper case "T".
reSoap              = '^' + reYear4 + '-' + reMonthLz + '-' + reDayLz + 'T'    + reHour24Lz + ':' + reMinuteLz + ':' + reSecondLz + reFrac + reTzCorrection + '?';
reXML_RPC           = '^' + reYear4       + reMonthLz + reDayLz + 'T'    + reHour24    + ':' + reMinuteLz + ':' + reSecondLz;
reXML_RPCNoColon    = '^' + reYear4       + reMonthLz + reDayLz + "[Tt]" + reHour24          + reMinuteLz + reSecondLz;
reWDDX              = '^' + reYear4 + '-' + reMonth   + '-' + reDay   + 'T'    + reHour24    + ':' + reMinute   + ':' + reSecond;
reEXIF              = '^' + reYear4 + ':' + reMonthLz + ':' + reDayLz + ' '    + reHour24Lz + ':' + reMinuteLz + ':' + reSecondLz;

rePgYearDotDay      = '^' + reYear4 + "\\.?" + reDayOfYear;
rePgTextShort       = "^(" + reMonthAbbr + ")-" + reDayLz + '-' + reYear;
rePgTextReverse     = '^' + "(\\d{3,4}|[4-9]\\d|3[2-9])"/*reYear*/ + "-(" + reMonthAbbr + ")-" + reDayLz;
reMssqlTime         = '^' + reHour12 + ":" + reMinuteLz + ":" + reSecondLz + "[:.]([0-9]+)" + reMeridian;
reISOWeekday        = '^' + reYear4 + "-?W" + reWeekOfYear + "-?([0-7])";
reISOWeek           = '^' + reYear4 + "-?W" + reWeekOfYear;

reFirstOrLastDay    = "^(first|last) day of";
reBackOrFrontOf    = "^(back|front) of " + reHour24 + reSpaceOpt + reMeridian + '?';
reYesterday        = "^yesterday";
reNow              = "^now";
reNoon             = "^noon";
reMidnightOrToday  = "^(midnight|today)";
reTomorrow         = "^tomorrow";

/* Common Log Format: 10/Oct/2000:13:55:36 -0700 */
reCLF               = '^' + reDay + "/(" + reMonthAbbr + ")/" + reYear4 + ':' + reHour24Lz + ':' + reMinuteLz + ':' + reSecondLz + reSpace + reTzCorrection;

/* Timestamp format: @1126396800 */
reTimestamp        = "^@(-?\\d+)";
reTimestampMs      = "^@(-?\\d+)\\.(\\d{0,6})"; // timestamp microsec

/* To fix some ambiguities */
reDateShortWithTimeShort12  = '^' + reDateNoYear + reTimeShort12.substring(1);
reDateShortWithTimeLong12   = '^' + reDateNoYear + reTimeLong12.substring(1);
reDateShortWithTimeShort    = '^' + reDateNoYear + reTimeShort24.substring(1);
reDateShortWithTimeLong     = '^' + reDateNoYear + reTimeLong24.substring(1);
reDateShortWithTimeLongTz   = '^' + reDateNoYear + reISO8601NormTz.substring(1);

/* Relative regexps */
reRelTextNumber = "first|second|third|fourth|fifth|sixth|seventh|eighth?|ninth|tenth|eleventh|twelfth";
reRelTextText   = "next|last|previous|this";
reRelTextUnit   = "(?:msec|millisecond|µsec|microsecond|usec|second|sec|minute|min|hour|day|fortnight|forthnight|month|year)s?|(?:ms)|(?:µs)|weeks|" + reDayText;

reRelNumber         = "([+-]*)[ \\t]*([0-9]{1,13})";
reRelative          = '^' + reRelNumber + reSpaceOpt + '(' + reRelTextUnit + "|week)";
reRelativeText      = "^(" + reRelTextNumber + '|' + reRelTextText + ')' + reSpace + '(' + reRelTextUnit + ')';
reRelativeTextWeek  = "^(" + reRelTextText + ')' + reSpace + "(week)";

reWeekdayOf = "^(" + reRelTextNumber + '|' + reRelTextText + ')' + reSpace + '(' + reDayFull + '|' + reDayAbbr + ')' + reSpace + "of";
```

## 平台支持

已经提供了 darwin_amd64、darwin_arm64、linux_amd64 的静态链接，如果需要在其他平台使用，自行编译 timelib 并添加 cgo 链接配置

```go
// strtotime.go

package timelib
/*
#include <stdlib.h>
#include "lib/timelib.h"
#include "strtotime.h"
#cgo darwin,arm64 LDFLAGS: -L./lib -ltimelib_darwin_arm64
#cgo darwin,amd64 LDFLAGS: -L./lib -ltimelib_darwin_amd64
#cgo linux,amd64 LDFLAGS: -L./lib -ltimelib_linux_amd64
*/
import "C"

// ...
```