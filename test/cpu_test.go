package test

import (
	"testing"
)

func TestCompile(t *testing.T) {
	compileElf(AsmDir, "addi")
}

func TestObjcopyBin(t *testing.T) {
	compileElf(AsmDir, "addi")
	objcopyBin(AsmDir, "addi")
}

func TestAddi(t *testing.T) {
	cpu := createCPU(AsmDir, "addi")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	x29 := -1
	x30 := -2
	assertEq(t, uint64(x29), cpu.Regs[29])
	assertEq(t, uint64(x30), cpu.Regs[30])
}

func TestAdd(t *testing.T) {
	cpu := createCPU(AsmDir, "add")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	x6 := 1
	x7 := -2
	x5 := x6 + x7
	assertEq(t, uint64(x6), cpu.Regs[6])
	assertEq(t, uint64(x7), cpu.Regs[7])
	assertEq(t, uint64(x5), cpu.Regs[5])
}

func TestAnd(t *testing.T) {
	cpu := createCPU(AsmDir, "and")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	x6 := 0x10
	x7 := 0x11
	x5 := x6 & x7
	assertEq(t, uint64(x6), cpu.Regs[6])
	assertEq(t, uint64(x7), cpu.Regs[7])
	assertEq(t, uint64(x5), cpu.Regs[5])
}

func TestAndi(t *testing.T) {
	cpu := createCPU(AsmDir, "andi")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	x6 := 0x91
	x5 := x6 & 0x01
	assertEq(t, uint64(x6), cpu.Regs[6])
	assertEq(t, uint64(x5), cpu.Regs[5])
}

func TestAuipc(t *testing.T) {
	cpu := createCPU(AsmDir, "auipc")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	x5 := cpu.Pc + 0x123<<12 - 8
	x6 := cpu.Pc - 4
	assertEq(t, uint64(x6), cpu.Regs[6])
	assertEq(t, uint64(x5), cpu.Regs[5])
}

func TestBne(t *testing.T) {
	cpu := createCPU(AsmDir, "bne")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	x5 := 6
	x6 := 6
	assertEq(t, uint64(x6), cpu.Regs[6])
	assertEq(t, uint64(x5), cpu.Regs[5])
}

func TestJalr(t *testing.T) {
	cpu := createCPU(CDir, "jalr")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	// fixme: find the bug of the command `jalr`
}

func TestFib(t *testing.T) {
	cpu := createCPU(AsmDir, "fib")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	// fixme: panic: invalid memory address: fffffffffffff810
}
