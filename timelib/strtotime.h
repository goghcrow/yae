#include "lib/timelib.h"
#include "lib/timelib_private.h"

long long strtotime(char *times, int len, long long preset_ts, timelib_tzinfo *tzi);