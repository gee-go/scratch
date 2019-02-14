package trie

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestTrie(t *testing.T) {
	var tr trie

	tr.set([]byte{'a', 'b'}, 1)
	tr.set([]byte{'a'}, 2)

	tr.Each(func(k []byte, v Value) {
		fmt.Printf("%s: %v\n", k, v)
	})

	f, err := os.Open("/usr/share/dict/words")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		tr.set(scanner.Bytes(), Value(i))
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	t.Log(int(ms.HeapAlloc))

	tr.Each(func(k []byte, v Value) {
		if v < 100 {
			fmt.Printf("%s: %v\n", k, v)
		}

	})
}
