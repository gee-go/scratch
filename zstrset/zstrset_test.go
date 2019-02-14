package zstrset

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	tags := []string{"abc:1", "abcd:2", "a:3"}
	z := makezstrset(tags)
	z.each(func(g, v []byte) {
		fmt.Printf("%s %s\n", g, v)
	})
	t.Fail()
}
