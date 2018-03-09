package cuckoo

import (
	"github.com/gee-go/cuckoo/hash2"
	"github.com/k0kubun/pp"
)

type bucket struct {
	k Key
	v Value

	// A key could be zero.
	// TODO: might be easier to have another index for this, assuming 0 will be a rare key.
	dirty bool // TODO remove
}

func nextPow2(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++

	return v
}

func nextPow64(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

const (
	maxLoop  = 32
	nbuckets = 2
)

type Bucket struct {
	mask uint32
	len  int
	data []bucket
}

type Cuckoo struct {
	t    [nbuckets][]bucket
	mask uint32
	n    uint32

	l       [nbuckets]int
	r       int
	seed    uint64
	scratch map[Key]Value
}

func (m *Cuckoo) Len() (l int) {
	for _, bl := range m.l {
		l += bl
	}
	return
}

func (m *Cuckoo) Cap() (l int) {
	for _, buckets := range m.t {
		l += len(buckets)
	}
	return
}

func NewCuckoo(n int) *Cuckoo {
	m := &Cuckoo{n: nextPow2(uint32(n)), scratch: make(map[Key]Value, 128)}
	m.rehash()
	return m
}

func (m *Cuckoo) rehash() {
	m.seed++
	m.r++
	if m.n == 0 {
		m.n = 1
	}

	for ti, tl := range m.l {
		if float64(tl)/float64(len(m.t[ti])) > 0.5 {
			m.n *= 2
			break
		}
	}
	// if float64(m.l)/float64(m.n) > 0.5 {
	// m.n *= 2
	// }

	bn := nextPow2(m.n / nbuckets)
	if bn < 2 {
		bn = 2
	}
	m.mask = bn - 1
	var old [nbuckets][]bucket
	for i := 0; i < nbuckets; i++ {
		old[i] = m.t[i]
		m.t[i] = make([]bucket, bn)
	}

	// speed up rehash by storing full hash?
	for _, buckets := range old {
		for _, b := range buckets {
			if b.dirty {
				m.insert(b, false)
			}
		}
	}

	for k, v := range m.scratch {
		m.insert(bucket{k: k, v: v, dirty: true}, false)
	}

	for k := range m.scratch {
		delete(m.scratch, k)
	}

}

func (m *Cuckoo) log(args ...interface{}) {
	pp.Println(args...)
}

func (m *Cuckoo) Get(k Key) (Value, bool) {
	h0, h1 := m.hash(k)

	b := m.t[0][h0]
	if b.k == k && b.dirty {
		return b.v, true
	}

	b = m.t[1][h1]
	if b.k == k && b.dirty {
		return b.v, true
	}

	v, ok := m.scratch[k]

	return v, ok
}

// func (m *Cuckoo) Get(k Key) (Value, bool) {
// 	var i int
// 	// we want the compiler to inline this function
// 	// so we use a goto
// 	// XXX: have test that makes sure this is inlined.
// loop:
// 	b0 := m.table(k, uint64(i))
// 	if b0.k == k && b0.dirty {
// 		return b0.v, true
// 	}
// 	i++

// 	if i < len(m.t) {
// 		goto loop
// 	}
// 	return 0, false
// }

func (m *Cuckoo) Insert(k Key, v Value) {
	b := bucket{k: k, v: v, dirty: true}

	m.insert(b, true)
}

// func (m *Cuckoo) table(k Key, i uint64) bucket {
// 	// keep separate so it inlines
// 	return m.t[i][hash(k, i)&m.mask]
// }

func (m *Cuckoo) insert(x bucket, first bool) {
	// lookup hash
	h0, h1 := m.hash(x.k)
	b := m.t[0][h0]
	if b.k == x.k && b.dirty {
		return
	}

	if !b.dirty {
		m.l[0]++
		m.t[0][h0] = x
		return
	}

	b = m.t[1][h1]
	if b.k == x.k && b.dirty {
		return
	}
	if !b.dirty {
		m.l[1]++
		m.t[1][h1] = x
		return
	}

	for i := 0; i < maxLoop; i++ {
		if i > 0 {
			h0, _ = m.hash(x.k)
		}

		p0 := m.t[0][h0]
		m.t[0][h0] = x
		// first slot open
		if !p0.dirty {
			m.l[0]++
			return
		}

		_, h1 = m.hash(p0.k)
		p1 := m.t[1][h1]
		m.t[1][h1] = p0
		// first slot open
		if !p1.dirty {
			m.l[1]++
			return
		}

		x = p1
	}

	// if len(m.scratch) < 1024 {
	// 	m.scratch[x.k] = x.v
	// 	return
	// }

	// for ti, t := range m.t {
	// 	c := 0
	// 	for _, b := range t {
	// 		if b.dirty {
	// 			c++
	// 		}
	// 	}

	// 	// pp.Println(ti, c, len(t), float64(c)/float64(len(t)))
	// 	// pp.Println(ti, m.l[0], m.l[1])
	// }
	// fmt.Println()

	m.rehash()

	// fmt.Println(m.l, m.r, m.n, float64(m.l)/float64(m.n), len(m.scratch))
	m.insert(x, false)
}

func (m *Cuckoo) hash(k Key) (uint32, uint32) {
	h := hash2.Int64Hash(uint64(k), m.seed)
	return uint32(h>>32) & m.mask, uint32(((h << 32) >> 32)) & m.mask
}
