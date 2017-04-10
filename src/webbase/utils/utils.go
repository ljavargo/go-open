package utils

import (
	"strconv"
	"strings"
)

func ParseImageSize(size string) int {
	size = strings.ToUpper(size)
	s := strings.TrimSuffix(size, "G")
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i
}
