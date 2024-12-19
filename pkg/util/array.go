package util

func HasTagAtOffset(data, tag []byte, dataOffset int) bool {
	for i, d := range tag {
		if data[dataOffset+i] != d {
			return false
		}
	}
	return true
}
