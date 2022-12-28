package main

import (
	"fmt"
	"goemu/config"
	"goemu/runtime"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("usage: goemu <filepath>")
		return
	}

	code, err := os.ReadFile(args[1])
	if err != nil {
		panic(err)
	}

	cpu := runtime.NewCPU(code)

	for cpu.Pc < config.KernelEnd {
		inst, err := cpu.Fetch()
		if err != nil {
			panic(err)
		}
		if err = cpu.Execute(inst); err != nil {
			panic(err)
		}
	}
}
