package runtime

import (
	"encoding/binary"
	"fmt"
	"goemu/config"
	"io"
)

type CPU struct {
	Regs [32]uint64
	Pc   uint64
	Mem  Memory
	Csr  CSR
	Level
}

func NewCPU(code []uint8) *CPU {
	regs := [32]uint64{}
	regs[2] = config.KernelEnd
	return &CPU{
		Regs: regs,
		Pc:   config.KernelBase,
		Mem:  code,
	}
}

func (cpu *CPU) Fetch() (inst InstructionFields, err error) {
	if cpu.Pc+4 > config.KernelBase+uint64(len(cpu.Mem)) {
		err = io.EOF
		return
	}
	if err = cpu.Mem.Check(cpu.Pc, 4); err != nil {
		return
	}
	index := cpu.Pc - config.KernelBase
	bits := binary.LittleEndian.Uint32(cpu.Mem[index : index+4])
	inst = ParseInstruction(bits)
	return
}

func (cpu *CPU) UpdatePC(nextPc *uint64) {
	cpu.Pc = *nextPc
}

// Execute function is part of the CPU struct and is responsible for executing a given instruction
func (cpu *CPU) Execute(inst InstructionFields) error {
	nextPc := cpu.Pc + 4 // add 4 by default
	defer cpu.UpdatePC(&nextPc)

	switch inst.Opcode {
	case 0b0000000: // nop
	case 0b0010011:
		shamt := uint32(inst.Rs2)
		switch inst.Funct3 {
		case 0b000: // addi
			cpu.Regs[inst.Rd] = cpu.Regs[inst.Rs1] + inst.ImmI
		case 0b001: // slli
			cpu.Regs[inst.Rd] = cpu.Regs[inst.Rd] << uint64(shamt)
		case 0b010: // slti
			if int64(cpu.Regs[inst.Rs1]) < int64(inst.ImmI) {
				cpu.Regs[inst.Rd] = 1
			} else {
				cpu.Regs[inst.Rd] = 0
			}
		case 0b011: //sltiu
			if cpu.Regs[inst.Rs1] < inst.ImmI {
				cpu.Regs[inst.Rd] = 1
			} else {
				cpu.Regs[inst.Rd] = 0
			}
		case 0b100: // xori
			cpu.Regs[inst.Rd] = cpu.Regs[inst.Rs1] ^ inst.ImmI
		case 0b101: // srli or srai
			switch inst.Funct7 {
			case 0b0000000: // srli
				cpu.Regs[inst.Rd] = cpu.Regs[inst.Rs1] >> shamt
			case 0b0100000: // srai
				cpu.Regs[inst.Rd] = uint64(int64(cpu.Regs[inst.Rs1]) >> shamt)
			default:
				return NewIllegalInstErr(inst)
			}
		case 0b110: // ori
			cpu.Regs[inst.Rd] = cpu.Regs[inst.Rs1] | inst.ImmI
		case 0b111: // andi
			cpu.Regs[inst.Rd] = cpu.Regs[inst.Rs1] & inst.ImmI
		}
	default:
		return NewIllegalInstErr(inst)
	}
	return nil
}

func NewIllegalInstErr(inst InstructionFields) error {
	return fmt.Errorf("unknown instruction format: %x", inst.Bits)
}
