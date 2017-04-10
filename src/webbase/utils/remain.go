package utils

import (
	"time"
)

func CaclRemainMonth(expire int64) float64 {
	now := time.Now().Unix()
	if now >= expire {
		return 0
	}

	poor := expire - now
	days := float64(poor) / 86400
	months := days / 30
	return months
}

func CaclRemainSecond(expire int64) uint64 {
	now := time.Now().Unix()
	if now >= expire {
		return 0
	}

	poor := expire - now
	return uint64(poor)
}
