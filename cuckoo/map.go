package cuckoo

import "github.com/gee-go/cuckoo/hash2"

type Key uint64

func (k Key) Hash(seed uint64) (uint32, uint32) {
	h := hash2.Int64Hash(uint64(k), seed)
	return uint32(h >> 32), uint32(((h << 32) >> 32))
}

type Value int

type Map interface {
	Get(Key) (Value, bool)
	Insert(Key, Value)
	Cap() int
	Len() int
}

type Std map[Key]Value

func (m Std) Get(k Key) (v Value, ok bool) {
	v, ok = m[k]
	return
}

func (m Std) Insert(k Key, v Value) {
	m[k] = v
}

func (m Std) Cap() int {
	return len(m) // TODO
}
func (m Std) Len() int {
	return len(m)
}
