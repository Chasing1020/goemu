package runtime

import "goemu/config"

type CPU struct {
	Regs [32]uint64
	Pc   uint64
	Mem  []uint8
}

func NewCPU(code []uint8) *CPU {
	regs := [32]uint64{}
	regs[2] = config.MemSize - 1
	return &CPU{
		Regs: regs,
		Pc:   0,
		Mem:  code,
	}
}
