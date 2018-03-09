package smhash

import (
	"testing"
)

func TestBanker(t *testing.T) {
	size := uint64(1) << 2
	x := Banker(1, size)
	for i := 0; i < 1; i++ {
		t.Logf("%d: %d %b", i, x, x)
		x = Banker(x, size)
	}
	// Gosper's hack
	// http://read.seas.harvard.edu/cs207/2012/?p=64

	// var s subset
	j := 0
	// var s subset
	s := Subset{N: 63, K: 5}
	for s.Next() {
		// s = s.Next(1)
		// t.Logf("%b", s.Val())
		j++
	}

	// i := 0
	// for s < (1 << 63) {
	// 	// t.Logf("%b", s)
	// 	s = s.Next(5)
	// 	i++
	// }

	t.Log(201376, j)

}

type Subset struct {
	v uint64
	N uint8
	K uint8
}

func (s *Subset) Val() uint64 {
	return s.v
}

func (s *Subset) Next() bool {
	if s.v == 0 {
		s.v = (1 << s.K) - 1
		return true
	}

	c := s.v & -s.v
	r := s.v + c

	s.v = (((r ^ s.v) >> 2) / c) | r
	return s.v < 1<<s.N
}
