package util

import (
	"time"
)

func TimeFromUnix(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, 0)
}
