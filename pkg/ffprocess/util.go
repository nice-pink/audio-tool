package ffprocess

import "strconv"

func GetFloatString(val float64) string {
	return strconv.FormatFloat(val, 'f', 3, 64)
}

func GetMilliSeconds(val float64) int {
	return int(val * 1000)
}
