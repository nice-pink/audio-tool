package ffprocess

import "strconv"

func GetFloatString(val float64) string {
	return strconv.FormatFloat(val, 'f', 4, 64)
}
