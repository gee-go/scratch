package bitset

import (
	"fmt"
	"testing"
)

func TestBitSet(t *testing.T) {
	var b dense

	// b.Set(0)

	// b.set(0)
	fmt.Println(b.w, b.wmin)
	// b.set(1)
	// b.set(0)
	// b.set(1)

	fmt.Println(b.w, b.wmin)
	// b.Set(256)
	// b.Set(129)
	// b.Set(112)

	// fmt.Println(b.data, b.offset)
	fmt.Println(b.test(0), b.test(2), b.test(64))

	t.Fatal(b.min())
}
