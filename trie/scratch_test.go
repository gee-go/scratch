package trie

// 1 item
// 8 + 1
// 8 + 8
type kset struct {
	data [16]*[16]*[16]*[16]uint8
}

// 2^16 / 16 = 4096
// 4096 / 16 = 256
// 256 / 16
