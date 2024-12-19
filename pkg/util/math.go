package util

func Unsynchsafe(in uint32) uint32 {
	var out uint32 = 0
	mask := uint32(0x7F000000)

	for {
		out >>= 1
		out |= in & mask
		mask >>= 8
		if mask == 0 {
			break
		}
	}

	return out
}

func Synchsafe(in uint32) uint32 {
	var out uint32 = 0
	mask := uint32(0x7F)

	for {
		out = in & ^mask
		out = out << 1
		out = out | in&mask
		mask = ((mask + 1) << 8) - 1
		in = out
		if mask^0x7FFFFFFF == 0 {
			break
		}
	}
	return out
}
