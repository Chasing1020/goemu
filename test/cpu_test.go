package test

import (
	"fmt"
	"goemu/runtime"
	"os"
	"testing"
)

func TestCPUFetch(t *testing.T) {
	filename := "a.bin"
	bits, err := os.ReadFile("./asm/" + filename)
	if err != nil {
		panic(err)
	}

	cpu := runtime.NewCPU(bits)
	fmt.Printf("%04x\n", cpu.Fetch())

	// Combine the uint8 values using bit shifting and bitwise OR
	result := (uint32(bits[3]) << 24) | (uint32(bits[2]) << 16) | (uint32(bits[1]) << 8) | uint32(bits[0])

	fmt.Printf("%04x", result)
}
