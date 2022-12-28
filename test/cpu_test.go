package test

import (
	"goemu/runtime"
	"runtime/debug"
	"testing"
)

func TestCompile(t *testing.T) {
	compileElf(AsmDir, "addi")
}

func TestObjcopyBin(t *testing.T) {
	compileElf(AsmDir, "addi")
	objcopyBin(AsmDir, "addi")
}

type Number interface {
	int8 | int16 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | uint
}

func assertEq[Number comparable](t *testing.T, expected, actual Number) {
	if expected != actual {
		t.Errorf("%s assertEq failed: expected %+v, got %+v\n%s", t.Name(), expected, actual, debug.Stack())
	}
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
	assertEq(t, x6, cpu.Regs[6])
	assertEq(t, x5, cpu.Regs[5])
}

func TestBlt(t *testing.T) {
	cpu := createCPU(AsmDir, "blt")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	assertEq(t, 1, cpu.Regs[8])
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
	cpu := createCPU(AsmDir, "jalr")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	assertEq(t, cpu.Regs[5], cpu.Regs[8])
	assertEq(t, 5, cpu.Regs[6])
}

func TestLa(t *testing.T) {
	cpu := createCPU(AsmDir, "la")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	assertEq(t, 114514, cpu.Regs[6])
}

func TestLb(t *testing.T) {
	cpu := createCPU(AsmDir, "lb")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	assertEq(t, 0x34, cpu.Regs[7])
}

func TestLh(t *testing.T) {
	cpu := createCPU(AsmDir, "lh")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	assertEq(t, 0x78, cpu.Regs[6])
	assertEq(t, 0x5678, cpu.Regs[7])
}

func TestLi(t *testing.T) {
	cpu := createCPU(AsmDir, "li")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	assertEq(t, 1<<(32<<(^uint32(0)>>63)-1)-1, cpu.Regs[7])
}

func TestLui(t *testing.T) {
	cpu := createCPU(AsmDir, "lui")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	assertEq(t, 0x12345<<12+0x678, cpu.Regs[5])
}

func TestNeg(t *testing.T) {
	cpu := createCPU(AsmDir, "neg")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	assertEq(t, cpu.Regs[7], cpu.Regs[5])
}

func TestNot(t *testing.T) {
	cpu := createCPU(AsmDir, "not")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	assertEq(t, cpu.Regs[7], cpu.Regs[5])
}

func TestCsr(t *testing.T) {
	cpu := createCPU(AsmDir, "csr")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	assertEq(t, 1, cpu.Csr[runtime.Mstatus])
	assertEq(t, 2, cpu.Csr[runtime.Mtvec])
	assertEq(t, 3, cpu.Csr[runtime.Mepc])
	assertEq(t, 0, cpu.Csr[runtime.Sstatus])
	assertEq(t, 5, cpu.Csr[runtime.Stvec])
	assertEq(t, 5, cpu.Csr[runtime.Sepc])
}

func TestSb(t *testing.T) {
	cpu := createCPU(AsmDir, "sb")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	// fixme: TestSb assertEq failed: expected 18, got 267461547
	assertEq(t, 0x12, cpu.Regs[6])
}

func TestFib(t *testing.T) {
	cpu := createCPU(AsmDir, "fib")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
	// fixme: panic: invalid memory address: fffffffffffff810
}
