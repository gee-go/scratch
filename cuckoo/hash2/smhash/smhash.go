package smhash

import (
	"math"
	"testing"
)

const hashSize = 32 + int(^uintptr(0)>>63<<5)

type Hasher interface {
	Bytes([]byte, uint64) uint64
	String(string, uint64) uint64
	Int64(uint64, uint64) uint64
}

type HashSet struct {
	m map[uint64]struct{} // set of hashes added
	n int                 // number of hashes added
	h Hasher
}

func NewHashSet(h Hasher) *HashSet {
	return &HashSet{
		m: make(map[uint64]struct{}),
		n: 0,
		h: h,
	}
}

func (s *HashSet) add(h uint64) {
	s.m[h] = struct{}{}
	s.n++
}
func (s *HashSet) addS(x string) {
	s.add(s.h.String(x, 0))
}
func (s *HashSet) addB(x []byte) {
	s.add(s.h.Bytes(x, 0))
}
func (s *HashSet) addS_seed(x string, seed uint64) {
	s.add(s.h.String(x, seed))
}
func (s *HashSet) check(t *testing.T) {
	const SLOP = 10.0
	collisions := s.n - len(s.m)
	t.Logf("%d/%d\n", len(s.m), s.n)
	pairs := int64(s.n) * int64(s.n-1) / 2
	expected := float64(pairs) / math.Pow(2.0, float64(hashSize))
	stddev := math.Sqrt(expected)
	if float64(collisions) > expected+SLOP*(3*stddev+1) {
		t.Errorf("unexpected number of collisions: got=%d mean=%f stddev=%f", collisions, expected, stddev)
	}
}
