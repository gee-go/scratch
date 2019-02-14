package bitset

import "math/bits"

const wsize = 32 << (^uint(0) >> 32 & 1)

type word uint

func (w *word) set(x uint16) {
	*w |= 1 << (x % wsize)
}

func (w word) test(x uint16) bool {
	return w&(1<<(x%wsize)) != 0
}

type dense struct {
	wmin int
	w    []word
}

func (d *dense) min() int {
	if len(d.w) == 0 {
		return -1
	}

	return wsize*d.wmin + bits.Len(uint(d.w[0])) - 1
}

func (d *dense) test(x uint16) bool {
	if len(d.w) == 0 {
		return false
	}

	slot := int(x / wsize)
	if slot < d.wmin {
		return false
	} else {
		slot -= d.wmin
	}

	if slot >= len(d.w) {
		return false
	}

	return d.w[slot].test(x)

}

func (d *dense) set(x uint16) {
	slot := int(x / wsize)

	if len(d.w) == 0 {
		d.wmin = slot
		slot = 0
	} else if slot < d.wmin {
		shift := d.wmin - slot
		d.wmin = slot
		slot = 0
		next := make([]word, len(d.w)+int(shift))
		copy(next[shift:], d.w)
		d.w = next
	} else {
		slot -= d.wmin
	}

	if slot >= len(d.w) {
		next := make([]word, slot+1)
		copy(next, d.w)
		d.w = next
	}

	d.w[slot].set(x)
}
