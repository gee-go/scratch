package wavelet

import (
	"fmt"
	"testing"
)

func Simple(in []int) []int {
	// |a| must be a power of 2

	out := make([]int, len(in))

	for i := 0; i < len(in)/2; i++ {
		a := in[i*2]
		b := in[i*2+1]

		s := (a + b) / 2
		d := a - s

		out[i] = s
		out[i+len(in)/2] = d
	}
	return out
	// fmt.Println(out)
}

func TestSimple(t *testing.T) {
	out := Simple([]int{56, 40, 8, 24, 48, 48, 40, 16})
	fmt.Println(out)
	Simple(out[:4])
	fmt.Println(out)

}
