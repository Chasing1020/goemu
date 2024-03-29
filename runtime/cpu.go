package runtime

import (
	"fmt"
	"goemu/config"
	"goemu/hw/uart"
	"goemu/util"
	"io"
	"strconv"
	"strings"
)

type CPU struct {
	Regs [32]uint64
	Pc   uint64
	Size uint64
	Bus  Bus
	Csr  CSR
	Level
}

func NewCPU(code []uint8) *CPU {
	regs := [32]uint64{}
	regs[2] = config.KernelEnd
	mem := make(Memory, config.MemSize)
	copy(mem, code)
	u := uart.NewUart()

	return &CPU{
		Regs: regs,
		Pc:   config.KernelBase,
		Size: uint64(len(code)),
		Bus: Bus{
			Mem:  &mem,
			Uart: u,
		},
		Level: Machine,
	}
}

// Run is a loop that fetches and executes instructions until an end-of-file error is encountered.
// It returns an error if there is a problem executing an instruction.
func (cpu *CPU) Run() error {
	for {
		inst, err := cpu.Fetch()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				panic(err)
			}
		}
		err = cpu.Execute(inst)
		if err != nil {
			util.ParseInst(inst)
			panic(err)
		}
	}
	return nil
}

func (cpu *CPU) Fetch() (inst uint64, err error) {
	if cpu.Pc < config.KernelBase || cpu.Pc >= config.KernelBase+cpu.Size {
		return 0, io.EOF
	}
	return cpu.Bus.Load(cpu.Pc, 4)
}

func (cpu *CPU) UpdatePC(nextPc *uint64) {
	cpu.Pc = *nextPc
}

// Execute function is part of the CPU struct and is responsible for executing a given instruction
func (cpu *CPU) Execute(inst uint64) error {
	nextPc := cpu.Pc + 4 // add 4 by default
	defer cpu.UpdatePC(&nextPc)

	opcode := uint8(inst & 0x0000007F)
	rs1 := uint8((inst & 0x000F8000) >> 15)
	rs2 := uint8((inst & 0x01F00000) >> 20)
	rd := uint8((inst & 0x00000F80) >> 7)
	funct3 := uint8((inst & 0x00007000) >> 12)
	funct7 := uint8((inst & 0xFE000000) >> 25)
	csrAddr := (inst & 0xFFF00000) >> 20

	immI := uint64(int32(inst&0xfff00000) >> 20)
	immS := (inst&0xFE000000)>>20 | (inst & 0x00000F80 >> 7)
	immB := uint64(int64(int32(inst&0x80000000)>>19)) | (inst & 0x80 << 4) | (inst >> 20 & 0x7E0) | (inst >> 7 & 0x1E)
	immJ := uint64((int32(uint64(inst)&0x80000000))>>11) | (uint64(inst) & 0xFF000) | ((inst >> 9) & 0x800) | ((inst >> 20) & 0x7FE)
	immU := inst & 0xFFFFF000

	switch opcode {
	case 0b0000011:
		addr := cpu.Regs[rs1] + immI
		switch funct3 {
		case 0b000: // lb
			val, err := cpu.Bus.Load(addr, 1)
			if err != nil {
				return err
			}
			cpu.Regs[rd] = uint64(int8(val))
		case 0b001: // lh
			val, err := cpu.Bus.Load(addr, 2)
			if err != nil {
				return err
			}
			cpu.Regs[rd] = uint64(int16(val))
		case 0b010: // lw
			val, err := cpu.Bus.Load(addr, 4)
			if err != nil {
				return err
			}
			cpu.Regs[rd] = uint64(int32(val))
		case 0b011: // ld
			val, err := cpu.Bus.Load(addr, 8)
			if err != nil {
				return err
			}
			cpu.Regs[rd] = val
		case 0b100: // lbu
			val, err := cpu.Bus.Load(addr, 1)
			if err != nil {
				return err
			}
			cpu.Regs[rd] = val
		case 0b101: // lhu
			val, err := cpu.Bus.Load(addr, 2)
			if err != nil {
				return err
			}
			cpu.Regs[rd] = val
		case 0b110: // lwu
			val, err := cpu.Bus.Load(addr, 4)
			if err != nil {
				return err
			}
			cpu.Regs[rd] = val
		default:
			return NewIllegalInstErr(inst)
		}
	case 0b0010011:
		switch funct3 {
		case 0b000: // addi
			cpu.Regs[rd] = cpu.Regs[rs1] + immI
		case 0b001: // slli
			cpu.Regs[rd] = cpu.Regs[rd] << uint64(rs2)
		case 0b010: // slti
			if int64(cpu.Regs[rs1]) < int64(immI) {
				cpu.Regs[rd] = 1
			} else {
				cpu.Regs[rd] = 0
			}
		case 0b011: //sltiu
			if cpu.Regs[rs1] < immI {
				cpu.Regs[rd] = 1
			} else {
				cpu.Regs[rd] = 0
			}
		case 0b100: // xori
			cpu.Regs[rd] = cpu.Regs[rs1] ^ immI
		case 0b101: // srli or srai
			switch funct7 {
			case 0b0000000: // srli
				cpu.Regs[rd] = cpu.Regs[rs1] >> rs2
			case 0b0100000: // srai
				cpu.Regs[rd] = uint64(int64(cpu.Regs[rs1]) >> rs2)
			default:
				return NewIllegalInstErr(inst)
			}
		case 0b110: // ori
			cpu.Regs[rd] = cpu.Regs[rs1] | immI
		case 0b111: // andi
			cpu.Regs[rd] = cpu.Regs[rs1] & immI
		}
	case 0b0010111: // auipc
		cpu.Regs[rd] = cpu.Pc + immU
	case 0b0011011:
		switch funct3 {
		case 0b000: // addiw
			cpu.Regs[rd] = uint64(int32(cpu.Regs[rs1] + immI))
		case 0b001: // slliw
			cpu.Regs[rd] = uint64(int32(cpu.Regs[rs1] << rs2))
		case 0b101:
			switch funct7 {
			case 0b0000000: // srliw
				cpu.Regs[rd] = uint64(int32(uint32(cpu.Regs[rs1]) >> rs2))
			case 0b0100000: // sraiw
				cpu.Regs[rd] = uint64(int32(cpu.Regs[rs1]) >> rs2)
			default:
				return NewIllegalInstErr(inst)
			}
		default:
			return NewIllegalInstErr(inst)
		}
	case 0b0100011:
		addr := cpu.Regs[rs1] + immS
		switch funct3 {
		case 0b000: // sb
			err := cpu.Bus.Store(addr, 1, cpu.Regs[rs2])
			if err != nil {
				return err
			}
		case 0b001: // sh
			err := cpu.Bus.Store(addr, 2, cpu.Regs[rs2])
			if err != nil {
				return err
			}
		case 0b010: // sw
			err := cpu.Bus.Store(addr, 4, cpu.Regs[rs2])
			if err != nil {
				return err
			}
		case 0b011: // sd
			err := cpu.Bus.Store(addr, 8, cpu.Regs[rs2])
			if err != nil {
				return err
			}
		default:
			return NewIllegalInstErr(inst)
		}
	case 0b0110011:
		switch funct3 {
		case 0b000:
			switch funct7 {
			case 0b0000000: // add
				cpu.Regs[rd] = cpu.Regs[rs1] + cpu.Regs[rs2]
			case 0b0000001: // mul
				cpu.Regs[rd] = cpu.Regs[rs1] * cpu.Regs[rs2]
			case 0b0100000: // sub
				cpu.Regs[rd] = cpu.Regs[rs1] - cpu.Regs[rs2]
			default:
				return NewIllegalInstErr(inst)
			}
		case 0b001: // sll
			cpu.Regs[rd] = cpu.Regs[rs1] << rs2
		case 0b010: // slt
			if int64(cpu.Regs[rs1]) < int64(cpu.Regs[rs2]) {
				cpu.Regs[rd] = 1
			} else {
				cpu.Regs[rd] = 0
			}
		case 0b011: // sltu
			if cpu.Regs[rs1] < cpu.Regs[rs2] {
				cpu.Regs[rd] = 1
			} else {
				cpu.Regs[rd] = 0
			}
		case 0b100: // xor
			cpu.Regs[rd] = cpu.Regs[rs1] ^ cpu.Regs[rs2]
		case 0b101:
			switch funct7 {
			case 0b0000000: // srl
				cpu.Regs[rd] = cpu.Regs[rs1] >> rs2
			case 0b0100000: // sra
				cpu.Regs[rd] = uint64(int64(cpu.Regs[rs1]) >> rs2)
			default:
				return NewIllegalInstErr(inst)
			}
		case 0b110: // or
			cpu.Regs[rd] = cpu.Regs[rs1] | cpu.Regs[rs2]
		case 0b111: // and
			cpu.Regs[rd] = cpu.Regs[rs1] & cpu.Regs[rs2]
		default:
			return NewIllegalInstErr(inst)
		}
	case 0b0110111: // lui
		cpu.Regs[rd] = immU
	case 0b0111011:
		switch funct3 {
		case 0b000:
			switch funct7 {
			case 0b0000000: // addw
				cpu.Regs[rd] = uint64(int64(int32(cpu.Regs[rs1] + cpu.Regs[rs2])))
			case 0b0100000: // subw
				cpu.Regs[rd] = uint64(int32(cpu.Regs[rs1] - cpu.Regs[rs2]))
			default:
				return NewIllegalInstErr(inst)
			}
		case 0b001: // sllw
			cpu.Regs[rd] = uint64(int32(uint32(cpu.Regs[rs1]) << rs2))
		case 0b101:
			switch funct7 {
			case 0b0000000: // srlw
				cpu.Regs[rd] = uint64(int32(uint32(cpu.Regs[rs1]) >> rs2))
			case 0b0000001: // divu
				if cpu.Regs[rs2] == 0 { // divisor equals zero
					cpu.Regs[rd] = 0xFFFFFFFFFFFFFFFF
				} else {
					cpu.Regs[rd] = cpu.Regs[rs1] / cpu.Regs[rs2]
				}
			case 0b0100000: // sraw
				cpu.Regs[rd] = uint64(int32(cpu.Regs[rs1]) >> int32(rs2))
			default:
				return NewIllegalInstErr(inst)
			}
		case 0b111:
			switch funct7 {
			case 0b0000001: // remuw
				if cpu.Regs[rs2] == 0 { // divisor equals zero
					cpu.Regs[rd] = cpu.Regs[rs1]
				} else {
					cpu.Regs[rd] = uint64(int32(cpu.Regs[rs1] % cpu.Regs[rs2]))
				}
			default:
				return NewIllegalInstErr(inst)
			}
		default:
			return NewIllegalInstErr(inst)
		}
	case 0b1100011:
		switch funct3 {
		case 0b000: // beq
			if cpu.Regs[rs1] == cpu.Regs[rs2] {
				nextPc = cpu.Pc + immB
			}
		case 0b001: // bne
			if cpu.Regs[rs1] != cpu.Regs[rs2] {
				nextPc = cpu.Pc + immB
			}
		case 0b100: // blt
			if int64(cpu.Regs[rs1]) < int64(cpu.Regs[rs2]) {
				nextPc = cpu.Pc + immB
			}
		case 0b101: // bge
			if int64(cpu.Regs[rs1]) >= int64(cpu.Regs[rs2]) {
				nextPc = cpu.Pc + immB
			}
		case 0b110: // bltu
			if cpu.Regs[rs1] < cpu.Regs[rs2] {
				nextPc = cpu.Pc + immB
			}
		case 0b111: // bgeu
			if cpu.Regs[rs1] >= cpu.Regs[rs2] {
				nextPc = cpu.Pc + immB
			}
		default:
			return NewIllegalInstErr(inst)
		}
	case 0b1100111: // jalr
		t := cpu.Pc + 4
		imm := uint64(int32(inst&0xFFF00000) >> 20)
		nextPc = (cpu.Regs[rs1] + imm) & ^(uint64(1))
		cpu.Regs[rd] = t
	case 0b1101111: // jal
		cpu.Regs[rd] = cpu.Pc + 4
		nextPc = cpu.Pc + immJ
	case 0b1110011:
		switch funct3 {
		case 0b000:
			switch funct7 {
			case 0b0001000: // sret
				sstatus, err := cpu.Csr.Load(Sstatus)
				if err != nil {
					return err
				}
				cpu.Level = Level((sstatus & SppMask) >> 8)
				spie := (sstatus & SpieMask) >> 5
				sstatus = (sstatus & ^uint64(SieMask)) | (spie << 1)
				sstatus |= SpieMask
				sstatus &= ^uint64(SppMask)
				err = cpu.Csr.Store(Sstatus, sstatus)
				if err != nil {
					return err
				}
				sepc, err := cpu.Csr.Load(Sepc)
				if err != nil {
					return err
				}
				nextPc = sepc &^ uint64(0b11)
			case 0b0011000: // mret
				mstatus, err := cpu.Csr.Load(Mstatus)
				if err != nil {
					return err
				}
				cpu.Level = Level((mstatus & MppMask) >> 11)
				mpie := (mstatus & MpieMask) >> 7
				mstatus = (mstatus & ^uint64(MieMask)) | (mpie << 3)
				mstatus |= MpieMask
				mstatus &= ^uint64(MppMask)
				mstatus &= ^uint64(MprvMask)
				err = cpu.Csr.Store(Mstatus, mstatus)
				if err != nil {
					return err
				}
				mepc, err := cpu.Csr.Load(Mepc)
				if err != nil {
					return err
				}
				nextPc = mepc & ^uint64(0b11)
			case 0b0001001: //sfence.vma
				return nil
			default:
				return NewIllegalInstErr(inst)
			}
		case 0b001: // csrrw
			data, err := cpu.Csr.Load(csrAddr)
			if err != nil {
				return err
			}
			if err = cpu.Csr.Store(csrAddr, cpu.Regs[rs1]); err != nil {
				return err
			}
			cpu.Regs[rd] = data
		case 0b010: // csrrs
			data, err := cpu.Csr.Load(csrAddr)
			if err != nil {
				return err
			}
			if err = cpu.Csr.Store(csrAddr, data|cpu.Regs[rs1]); err != nil {
				return err
			}
			cpu.Regs[rd] = data
		case 0b011: // csrrc
			data, err := cpu.Csr.Load(csrAddr)
			if err != nil {
				return err
			}
			if err = cpu.Csr.Store(csrAddr, data&(^cpu.Regs[rs1])); err != nil {
				return err
			}
			cpu.Regs[rd] = data
		case 0b101: // csrrwi
			data, err := cpu.Csr.Load(csrAddr)
			if err != nil {
				return err
			}
			cpu.Regs[rd] = data
			if err = cpu.Csr.Store(csrAddr, uint64(rs1)); err != nil {
				return err
			}
		case 0b110: // csrrsi
			data, err := cpu.Csr.Load(csrAddr)
			if err != nil {
				return err
			}
			if err = cpu.Csr.Store(csrAddr, data|uint64(rs1)); err != nil {
				return err
			}
			cpu.Regs[rd] = data
		case 0b111: // csrrci
			data, err := cpu.Csr.Load(csrAddr)
			if err != nil {
				return err
			}
			if err = cpu.Csr.Store(csrAddr, data&(^uint64(rs1))); err != nil {
				return err
			}
			cpu.Regs[rd] = data
		default:
			return NewIllegalInstErr(inst)
		}
	default:
		return NewIllegalInstErr(inst)
	}
	return nil
}

var (
	AbiMap = [32]string{
		"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2",
		"s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
		"a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
		"s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
	}
)

// GetReg is a method of the CPU struct that takes a string input and returns an uint64 value and an error.
func (cpu *CPU) GetReg(s string) (data uint64, err error) {
	i := -1
	if s[0] == 'x' || s[0] == 'X' {
		i, err = strconv.Atoi(s[1:])
		if err != nil {
			return
		}
	} else {
		for k, v := range AbiMap {
			if strings.EqualFold(v, s) {
				i = k
			}
		}
	}
	if i < 0 || i >= 32 {
		return 0, fmt.Errorf("unknown register format: %s", s)
	}
	data = cpu.Regs[i]
	return
}

func (cpu *CPU) debug(inst uint32) {
	fmt.Printf("inst: %08x, Pc: %08x\n", inst, cpu.Pc)
	fmt.Printf("inst: %032b, Pc: %032b\n", inst, cpu.Pc)
}

func NewIllegalInstErr(inst uint64) error {
	return fmt.Errorf("unknown instruction format: %x", inst)
}
