package hash2

import (
	_ "runtime"
	"unsafe" // required to use //go:linkname
)

//go:noescape
//go:linkname useAeshash runtime.useAeshash
var useAeshash bool

//go:noescape
//go:linkname strhash runtime.strhash
func strhash(a unsafe.Pointer, h uintptr) uintptr

//go:noescape
//go:linkname aeshash runtime.aeshash
func aeshash(p unsafe.Pointer, h uintptr) uintptr

//go:noescape
//go:linkname aeshash64 runtime.aeshash64
func aeshash64(p unsafe.Pointer, h uintptr) uintptr

//go:noescape
//go:linkname int64Hash runtime.int64Hash
func int64Hash(i uint64, seed uintptr) uintptr

//go:noescape
//go:linkname int32Hash runtime.int32Hash
func int32Hash(i uint32, seed uintptr) uintptr

//go:noescape
//go:linkname stringHash runtime.stringHash
func stringHash(s string, seed uintptr) uintptr

//go:noescape
//go:linkname bytesHash runtime.bytesHash
func bytesHash(b []byte, seed uintptr) uintptr

//
func Int64Hash(i uint64, seed uint64) uint64 {
	return uint64(int64Hash(i, uintptr(seed)))
}

// func Int32Hash(i uint32, seed uint64) uint64 {
// 	return uint64(int32Hash(i, uintptr(seed)))
// }

// func StringHash(v string, seed uint64) uint64 {
// 	return uint64(stringHash(v, uintptr(seed)))
// }

// func BytesHash(v []byte, seed uint64) uint64 {
// 	return uint64(bytesHash(v, uintptr(seed)))
// }
