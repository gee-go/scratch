package cuckoo

import (
	"fmt"
	"testing"
)

func TestNNextPow64(t *testing.T) {
	for i := 0; i < 20; i++ {
		t.Log(nextPow64(uint64(i)))
	}

}

func TestCuckoo(t *testing.T) {
	m := NewCuckoo(8)
	n := 1000
	for i := 0; i < n; i++ {
		m.Insert(Key(i), Value(i))
		m.Insert(Key(i), Value(i))
	}

	for i := 0; i < n; i++ {
		v, ok := m.Get(Key(i))
		if !ok {
			t.Fatal(i, "not found")
		}

		if v != Value(i) {
			t.Fatal("got", v, "expected", i)
		}
	}

	t.Log(float64(m.Cap()), m.l, m.r)
}

const (
	benchMod = 1000000
)

func makeMap(t string, n int) Map {
	switch t {
	case "std":
		return make(Std, n)
	case "cuckoo":
		return NewCuckoo(n)
	}

	panic(t)
}

type benchConfig struct {
	sizeHint int
	nValues  int
}

func (b benchConfig) String() string {
	return fmt.Sprintf("v=%d,size=%d", b.nValues, b.sizeHint)
}

func benchInsert(b *testing.B, name string, c benchConfig) {
	b.Run(fmt.Sprintf("%s.insert.%s", name, c), func(b *testing.B) {
		m := makeMap(name, c.sizeHint)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			k := Key(i % c.nValues)
			m.Insert(k, Value(i))

			if v, ok := m.Get(k); !ok {
				b.Fatal(v)
			}
		}
	})

	b.Run(fmt.Sprintf("%s.grow.%d", name, c.nValues), func(b *testing.B) {
		m := makeMap(name, 1)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			k := Key(i % c.nValues)
			m.Insert(k, Value(i))

			if v, ok := m.Get(k); !ok {
				b.Fatal(v)
			}
		}
	})
}

func benchGet(b *testing.B, name string, c benchConfig) {
	b.Run(fmt.Sprintf("%s.get.%s", name, c), func(b *testing.B) {
		m := makeMap(name, c.sizeHint)
		for i := 0; i < c.nValues; i++ {
			m.Insert(Key(i), Value(i))
		}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			k := Key(i % c.nValues)

			if v, ok := m.Get(k); !ok {
				b.Fatal(v)
			}
		}
	})
}

func BenchmarkMap(b *testing.B) {
	for _, name := range []string{"cuckoo", "std"} {
		cfg := benchConfig{
			nValues: 10000000,
		}
		cfg.sizeHint = int(1.1 * float64(cfg.nValues))
		benchInsert(b, name, cfg)
		benchGet(b, name, cfg)
	}

	// benchmark(b, "cuckoo")
	// benchmark(b, "std")

}
