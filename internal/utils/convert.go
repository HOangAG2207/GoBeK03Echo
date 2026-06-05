package utils

import "time"

func Int64ToDuration(exp int64) time.Duration {
	if exp <= 0 {
		return 0 // no expiration
	}
	return time.Duration(exp) * time.Second
}
