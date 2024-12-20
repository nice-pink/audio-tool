package util

import (
	"strconv"
	"time"
)

func GetFilePath(baseFilePath string) string {
	if baseFilePath == "" {
		return ""
	}
	now := time.Now()
	return baseFilePath + "_" + strconv.FormatInt(now.Unix(), 10)
}
