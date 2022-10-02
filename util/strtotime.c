#include "lib/timelib.h"
#include "lib/timelib_private.h"

extern timelib_tzinfo *parse_tzfile(const char *formal_tzname, const timelib_tzdb *tzdb, int *dummy_error_code);

// 移植下 php 的 strtotime
long long strtotime(char *times, int len, long long preset_ts, timelib_tzinfo *tzi) {
	int parse_error, epoch_does_not_fit_in;
	timelib_error_container *error;
	long long ts;
	timelib_time *t, *now;

	if (len <= 0) {
		return 0;
	}
	if (!tzi) {
	    return 0;
	}

	now = timelib_time_ctor();
	now->tz_info = tzi;
	now->zone_type = TIMELIB_ZONETYPE_ID;
	timelib_unixtime2local(now, (timelib_sll) preset_ts);

	t = timelib_strtotime(times, len, &error, timelib_builtin_db(), parse_tzfile);
	parse_error = error->error_count;
	// 忽略错误
	timelib_error_container_dtor(error);
	if (parse_error) {
		timelib_time_dtor(now);
		timelib_time_dtor(t);
		return 0;
	}

	timelib_fill_holes(t, now, TIMELIB_NO_CLONE);
	timelib_update_ts(t, tzi);
	ts = timelib_date_to_int(t, &epoch_does_not_fit_in);

	timelib_time_dtor(now);
	timelib_time_dtor(t);

	if (epoch_does_not_fit_in) {
		// 忽略错误
		return 0;
	}

	return ts;
}