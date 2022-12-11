package test

import (
	"testing"
)

func TestCompile(t *testing.T) {
	compileElf("addi")
}

func TestObjcopyBin(t *testing.T) {
	compileElf("addi")
	objcopyBin("addi")
}

func TestAddi(t *testing.T) {
	cpu := createCPU("addi")
	if err := cpu.Run(); err != nil {
		panic(err)
	}
	x29 := -1
	x30 := -2
	assertEq(t, uint64(x29), cpu.Regs[29])
	assertEq(t, uint64(x30), cpu.Regs[30])
}

func TestAdd(t *testing.T) {
	cpu := createCPU("add")
	if err := cpu.Run(); err != nil {
		panic(err)
	}
	x6 := 1
	x7 := -2
	x5 := x6 + x7
	assertEq(t, uint64(x6), cpu.Regs[6])
	assertEq(t, uint64(x7), cpu.Regs[7])
	assertNeq(t, uint64(x5), cpu.Regs[5])
}
