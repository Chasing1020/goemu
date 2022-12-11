package test

import (
	"goemu/runtime"
	"os"
	"os/exec"
	"runtime/debug"
	"testing"
)

var (
	CrossCompile = "riscv64-unknown-elf-"
	Cflags       = []string{"-nostdlib", "-fno-builtin", "-march=rv64g", "-mabi=lp64", "-g", "-Wall", "-Ttext=0x80000000"}
	CC           = CrossCompile + "gcc"
	Objcopy      = CrossCompile + "objcopy"

	Pwd, _ = os.Getwd()
	AsmDir = Pwd + "/asm/"
)

func compileElf(asmName string) {
	cmd := exec.Command(CC, append(Cflags, "-o", asmName+".elf", asmName+".s")...)
	cmd.Dir = Pwd + "/asm"
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func objcopyBin(asmName string) {
	cmd := exec.Command(Objcopy, "-O", "binary", asmName+".elf", asmName+".bin")
	cmd.Dir = Pwd + "/asm"
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func createCPUFromCode(asmName, code string) *runtime.CPU {
	if err := os.WriteFile(AsmDir+asmName+".s", []byte(code), 0666); err != nil {
		panic(err)
	}
	compileElf(asmName)
	objcopyBin(asmName)
	bits, err := os.ReadFile(AsmDir + asmName + ".bin")
	if err != nil {
		panic(err)
	}
	return runtime.NewCPU(bits)
}

func createCPU(asmName string) *runtime.CPU {
	compileElf(asmName)
	objcopyBin(asmName)
	bits, err := os.ReadFile(AsmDir + asmName + ".bin")
	if err != nil {
		panic(err)
	}
	return runtime.NewCPU(bits)
}

type Number interface {
	int8 | int16 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | uint
}

func assertEq[Number comparable](t *testing.T, expected, actual Number) {
	if expected != actual {
		t.Errorf("%s assertEq failed: expected %+v, got %+v\n%s", t.Name(), expected, actual, debug.Stack())
	}
}

func assertNeq[Number comparable](t *testing.T, expected, actual Number) {
	if expected == actual {
		t.Errorf("%s assertNeq failed: expected %+v, got %+v\n%s", t.Name(), expected, actual, debug.Stack())
	}
}
