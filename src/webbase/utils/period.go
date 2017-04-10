package utils

import (
	"time"
)

func GetExpireTime() int64 {
	n := time.Now()
	return n.Unix()
}
