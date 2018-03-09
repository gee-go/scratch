package hash2

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/DataDog/mmh3"
	farm "github.com/dgryski/go-farm"
)

func ddmmh(k1, k2 uint64) (h1, h2 uint64) {
	k := make([]byte, 16)
	binary.LittleEndian.PutUint64(k[:], k1)
	binary.LittleEndian.PutUint64(k[8:], k2)

	h := mmh3.Hash128x64(k)

	return binary.LittleEndian.Uint64(h[:]), binary.LittleEndian.Uint64(h[8:])
}

func TestMurmur3(t *testing.T) {
	const (
		k1, k2 uint64 = 1, 2314
	)

	e1, e2 := ddmmh(k1, k2)
	a1, a2 := Hash128(k1, k2)

	if e1 != a1 || e2 != a2 {
		t.Fatalf("%d,%d != %d,%d", e1, e2, a1, a2)
	}
}

// func TestHash64(t *testing.T) {
// 	const n = 1 << 30
// 	pp.Println(int(exp(1<<30, 1<<30)))
// 	if testing.Short() {
// 		t.Skip("short")
// 	}

// 	// out := make(map[uint64]struct{}, n*2)

// 	var out [n]uint32

// 	collide := 0
// 	for i := uint64(0); i < n; i++ {
// 		v := Hash64(i) & (n - 1)
// 		if out[v] != 0 {
// 			collide++
// 		}

// 		out[v]++
// 	}

// 	t.Log(n, collide, float64(n)/float64(collide))
// }

func exp(n, d float64) float64 {
	return n - d + math.Pow((d-1)/d, n)

}

var guint64 uint64
var guintp uintptr
var guint32 uint32

func BenchmarkMurmur3(b *testing.B) {

	b.Run("uint64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			guint64 = Hash64(uint64(i))
		}
	})
}

func BenchmarkHash128(b *testing.B) {
	raw := make([]byte, 64)
	for i := range raw {
		raw[i] = 'a'
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// guint32 = mmh3.Hash32(raw)
		guint32 = farm.Fingerprint32(raw)
		// guint64 = farm.Hash64WithSeed(raw, 0)
	}

}

type benchCase struct {
	name string
	str  func(string, uint64) uint64
}

func benchString(b *testing.B, n int, bc benchCase) {
	b.Run(fmt.Sprintf("string-%d-%s", n, bc.name), func(b *testing.B) {
		raw := make([]byte, n)
		for i := range raw {
			raw[i] = 'a'
		}
		key := string(raw)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			guint64 = bc.str(key, 0)
		}
	})
}

// func BenchmarkRuntime(b *testing.B) {
// 	b.Run("uint64", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			guint64 = Int64Hash(uint64(i), 0)
// 		}
// 	})

// 	for _, bc := range []benchCase{
// 		{
// 			name: "runtime",
// 			str:  StringHash,
// 		},
// 		// {
// 		// 	name: "murmur3",
// 		// 	str:  mmhString,
// 		// },
// 	} {

// 		for i := 1; i <= 256; i *= 2 {
// 			benchString(b, i, bc)
// 		}
// 	}

// }

func TestHash(t *testing.T) {
	t.Log(useAeshash)
	for i := 0; i < 10; i++ {
		t.Logf("%d: %b", i, uint32(aeshash64(unsafe.Pointer(&i), 44)))
	}
}
