package runtime

import (
	"encoding/binary"
	"goemu/config"
)

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

func (cpu *CPU) Fetch() uint32 {
	defer func() { cpu.Pc += 4 }()
	return binary.LittleEndian.Uint32(cpu.Mem[cpu.Pc : cpu.Pc+4])
}

func (cpu *CPU) Decode() uint32 {
	panic("unimplemented")
}

func (cpu *CPU) Execute(bits uint32) uint32 {
	panic("unimplemented")
}
