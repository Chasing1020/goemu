package test

import (
	"fmt"
	"testing"
)

func TestCPUFetch(t *testing.T) {
	a := uint8(1)
	b := uint8(2)
	c := uint8(3)
	d := uint8(4)

	// Combine the uint8 values using bit shifting and bitwise OR
	result := (uint32(a) << 24) | (uint32(b) << 16) | (uint32(c) << 8) | uint32(d)

	fmt.Println(result) // Output: 16909060
}
