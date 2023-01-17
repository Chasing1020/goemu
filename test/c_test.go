package test

import (
	"goemu/runtime"
	"goemu/util"
	"os"
	"testing"
)

func newCRuntime(name string) *runtime.CPU {
	pwd, _ := os.Getwd()
	filepath := pwd + "/c/"
	util.Compile(filepath+name+".c", filepath+name+".elf")
	util.Objcopy(filepath+name+".elf", filepath+name+".bin")
	bits, err := os.ReadFile(filepath + name + ".bin")
	if err != nil {
		panic(err)
	}
	return runtime.NewCPU(bits)
}

func TestHello(t *testing.T) {
	cpu := newCRuntime("hello")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestInput(t *testing.T) {
	cpu := newCRuntime("input")
	if err := cpu.Run(); err != nil {
		t.Fatal(err)
	}
}
