package hash2

import (
	"math/bits"
	"reflect"
	"unsafe"
)

const (
	h32c1 uint32 = 0xcc9e2d51
	h32c2 uint32 = 0x1b873593
	h64c1 uint64 = 0x87c37b91114253d5
	h64c2 uint64 = 0x4cf5ad432745937f
)

func Mix32(h uint32) uint32 {
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16

	return h
}

func Mix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33

	return k
}

func Hash128(k1, k2 uint64) (h1, h2 uint64) {
	k1 *= h64c1
	k1 = bits.RotateLeft64(k1, 31)
	k1 *= h64c2
	h1 ^= k1
	h1 = bits.RotateLeft64(h1, 27)
	h1 += h2
	h1 = h1*5 + 0x52dce729
	k2 *= h64c2
	k2 = bits.RotateLeft64(k2, 33)
	k2 *= h64c1
	h2 ^= k2
	h2 = bits.RotateLeft64(h2, 31)
	h2 += h1
	h2 = h2*5 + 0x38495ab5

	// mix
	h1 ^= uint64(16)
	h2 ^= uint64(16)
	h1 += h2
	h2 += h1
	h1 = Mix64(h1)
	h2 = Mix64(h2)
	h1 += h2
	h2 += h1

	return
}

// Hash64 is a unholy combination of MurmurHash3_x86_32 and MurmurHash3_x64_128
// TODO: test that this is valid
func Hash64(k1 uint64) uint64 {
	var h1 uint64
	k1 *= h64c1
	k1 = bits.RotateLeft64(k1, 31)
	k1 *= h64c2
	h1 ^= k1
	h1 = bits.RotateLeft64(h1, 27)
	h1 = h1*5 + 0x52dce729

	// finalization
	h1 ^= 8
	return Mix64(h1)
}

func Hash64Seed(k1, seed uint64) uint64 {
	h1 := seed
	k1 *= h64c1
	k1 = bits.RotateLeft64(k1, 31)
	k1 *= h64c2
	h1 ^= k1
	h1 = bits.RotateLeft64(h1, 27)
	h1 = h1*5 + 0x52dce729

	// finalization
	h1 ^= 8
	return Mix64(h1)
}

// Hash128x64 is a version of MurmurHash which is designed to run only on
// little-endian processors.  It is considerably faster for those processors
// than Hash128.
func Hash128x64(key []byte, seed uint64) (uint64, uint64) {
	length := len(key)

	nblocks := length / 16
	var k1, k2 uint64
	h1, h2 := seed, seed
	h := *(*reflect.SliceHeader)(unsafe.Pointer(&key))
	h.Len = nblocks * 2
	b := *(*[]uint64)(unsafe.Pointer(&h))
	for i := 0; i < len(b); i += 2 {
		k1, k2 = b[i], b[i+1]
		k1 *= h64c1
		k1 = (k1 << 31) | (k1 >> (64 - 31))
		k1 *= h64c2
		h1 ^= k1
		h1 = (h1 << 27) | (h1 >> (64 - 27))
		h1 += h2
		h1 = h1*5 + 0x52dce729
		k2 *= h64c2
		k2 = (k2 << 33) | (k2 >> (64 - 33))
		k2 *= h64c1
		h2 ^= k2
		h2 = (h2 << 31) | (h2 >> (64 - 31))
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}
	h.Len = length

	k1, k2 = 0, 0
	tailIndex := nblocks * 16
	switch length & 15 {
	case 15:
		k2 ^= uint64(key[tailIndex+14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(key[tailIndex+13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(key[tailIndex+12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(key[tailIndex+11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(key[tailIndex+10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(key[tailIndex+9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(key[tailIndex+8])
		k2 *= h64c2
		k2 = (k2 << 33) | (k2 >> (64 - 33))
		k2 *= h64c1
		h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(key[tailIndex+7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(key[tailIndex+6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(key[tailIndex+5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(key[tailIndex+4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(key[tailIndex+3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(key[tailIndex+2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(key[tailIndex+1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(key[tailIndex])
		k1 *= h64c1
		k1 = (k1 << 31) | (k1 >> (64 - 31))
		k1 *= h64c2
		h1 ^= k1
	}
	h1 ^= uint64(length)
	h2 ^= uint64(length)
	h1 += h2
	h2 += h1
	h1 ^= h1 >> 33
	h1 *= 0xff51afd7ed558ccd
	h1 ^= h1 >> 33
	h1 *= 0xc4ceb9fe1a85ec53
	h1 ^= h1 >> 33
	h2 ^= h2 >> 33
	h2 *= 0xff51afd7ed558ccd
	h2 ^= h2 >> 33
	h2 *= 0xc4ceb9fe1a85ec53
	h2 ^= h2 >> 33
	h1 += h2
	h2 += h1

	return h1, h2
}
