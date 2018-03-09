package smhash

func Banker(x, highbit uint64) uint64 {
	y := x & -x
	c := y + x

	x = (((x ^ c) >> 2) / y) | c
	if x&highbit > 0 {
		x = ((x & (highbit - 1)) << 2) | 3
	}
	return x
}
