package test

import (
	"goemu/runtime"
	"os"
	"os/exec"
)

var (
	CrossCompile = "riscv64-unknown-elf-"
	Cflags       = []string{"-nostdlib", "-fno-builtin", "-march=rv64g", "-mabi=lp64", "-g", "-Wall", "-Ttext=0x80000000"}
	CC           = CrossCompile + "gcc"
	Objcopy      = CrossCompile + "objcopy"

	Pwd, _ = os.Getwd()
	AsmDir = Pwd + "/asm/"
	CDir   = Pwd + "/c/"
)

func compileElf(dir, name string) {
	if dir == AsmDir {
		cmd := exec.Command(CC, append(Cflags, "-o", name+".elf", name+".s")...)
		cmd.Dir = Pwd + "/asm"
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
	if dir == CDir {
		cmd := exec.Command(CC, append(Cflags, "-o", name+".elf", name+".c")...)
		cmd.Dir = Pwd + "/c"
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}

func objcopyBin(dir, name string) {
	cmd := exec.Command(Objcopy, "-O", "binary", name+".elf", name+".bin")
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func createCPU(dir, name string) *runtime.CPU {
	compileElf(dir, name)
	objcopyBin(dir, name)
	bits, err := os.ReadFile(dir + name + ".bin")
	if err != nil {
		panic(err)
	}
	return runtime.NewCPU(bits)
}
