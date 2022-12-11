package test

import (
	"fmt"
	"goemu/runtime"
	"io"
	"os"
	"os/exec"
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

func CompileElf(asmName string) {
	cmd := exec.Command(CC, append(Cflags, "-o", asmName+".elf", asmName+".s")...)
	cmd.Dir = Pwd + "/asm"
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func TestCompile(t *testing.T) {
	CompileElf("addi")
}

func ObjcopyBin(asmName string) {
	cmd := exec.Command(Objcopy, "-O", "binary", asmName+".elf", asmName+".bin")
	cmd.Dir = Pwd + "/asm"
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func TestObjcopyBin(t *testing.T) {
	CompileElf("addi")
	ObjcopyBin("addi")
}

func CreateCPU(asmName, code string) *runtime.CPU {
	if err := os.WriteFile(AsmDir+asmName+".s", []byte(code), 0666); err != nil {
		panic(err)
	}
	CompileElf(asmName)
	ObjcopyBin(asmName)
	bits, err := os.ReadFile(AsmDir + asmName + ".bin")
	if err != nil {
		panic(err)
	}
	return runtime.NewCPU(bits)
}

func TestAddi(t *testing.T) {
	asmName := "add"
	code := `.text
.global	_start

_start:
    addi x29, x0, 114
    addi x30, x0, 514
`
	cpu := CreateCPU(asmName, code)
	for {
		inst, err := cpu.Fetch()
		if err == io.EOF {
			break
		}
		err = cpu.Execute(inst)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("x29: %d\n", cpu.Regs[29])
	fmt.Printf("x30: %d\n", cpu.Regs[30])
}
