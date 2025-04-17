package util

import (
	"strconv"
	"time"
)

func CurrentTimestampMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func CurrentTimestampMillisString() string {
	return strconv.FormatInt(CurrentTimestampMillis(), 10)
}
