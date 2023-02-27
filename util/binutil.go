package util

import "os/exec"

var (
	//crossCompile = "riscv64-linux-gnu-"
	crossCompile = "riscv64-unknown-elf-"
	cflags       = []string{"-nostdlib", "-fno-builtin", "-march=rv64g", "-mabi=lp64", "-O1", "-Wall", "-Ttext=0x80000000"}
	cc           = crossCompile + "gcc"

	objcopy = "llvm-objcopy"
	objdump = "llvm-objdump"
)

func Compile(infile, outfile string) {
	cmd := exec.Command(cc, append(cflags, "-o", outfile, infile)...)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func Objcopy(infile, outfile string) {
	cmd := exec.Command(objcopy, "-j", ".text", "-O", "binary", infile, outfile)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func Objdump(infile string) {
	cmd := exec.Command(objdump, "-M", "no-aliases", "-d", infile)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
