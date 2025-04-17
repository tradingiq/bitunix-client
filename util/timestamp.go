package util

import (
	"time"
)

func CurrentTimestampMillis() int64 {
	return time.Now().UnixMilli()
}
