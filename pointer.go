package main

import (
	"fmt"
	"unsafe"
)

func main() {
	a, b, c := 0, 1, 2
	fmt.Printf("%d\n", uintptr(unsafe.Pointer(&a)))
	fmt.Printf("%d\n", uintptr(unsafe.Pointer(&b))-uintptr(unsafe.Pointer(&a)))
	fmt.Printf("%d\n", uintptr(unsafe.Pointer(&c)))

}
