package runtime

// Instruction includes fields for all possible instruction types
// can be a convenient way to represent RISC-V instructions
// not all fields will be valid for all instruction types
type Instruction struct {
	Opcode uint8
	Rs1    uint8
	Rs2    uint8
	Rd     uint8
	Func3  uint8
	Func7  uint8
	Csr    uint32
	ImmI   int32
	ImmS   int32
	ImmB   int32
	ImmJ   int32
	ImmU   int32
}

type InstructionType uint8

const (
	Unknown InstructionType = iota

	LUI
	AUIPC
	JAL
	JALR
	BEQ
	BNE
	BLT
	BGE
	BLTU
	BGEU
	LB
	LH
	LW
	LD
	LBU
	LHU
	LWU
	SB
	SH
	SW
	SD
	ADDI
	SLTI
	SLTIU
	XORI
	ORI
	ANDI
	ADDIW
	SLLIW
	SRLIW
	SRAIW
	SLLI
	SRLI
	SRAI
	ADD
	MUL
	SUB
	SLL
	SLT
	SLTU
	XOR
	SRL
	SRA
	OR
	AND
	FENCE
	FENCEI
	CSRRW
	CSRRS
	CSRRC
	CSRRWI
	CSRRSI
	CSRRCI
	ECALL
	EBREAK
	ERET
	WFI
)

// ParseInstruction extracts the instruction's opcode, register fields, and immediate values by masking and shifting the input bits.
func ParseInstruction(bits uint32) (inst Instruction) {
	inst.Opcode = uint8(bits & 0x0000007F)
	inst.Rs1 = uint8((bits & 0x000F8000) >> 15)
	inst.Rs2 = uint8((bits & 0x01F00000) >> 20)
	inst.Rd = uint8((bits & 0x00000F80) >> 7)
	inst.Func3 = uint8((bits & 0x00007000) >> 12)
	inst.Func7 = uint8((bits & 0xFE000000) >> 25)
	inst.Csr = (bits & 0xFFF00000) >> 20

	inst.ImmI = int32(bits&0xFFF00000) >> 20
	inst.ImmS = (int32(bits&0xFE000000) >> 20) | (int32(bits&0x00000F80) >> 7)
	inst.ImmB = (int32(bits&0x80000000) >> 19) | (int32(bits&0x00000080) << 4) | (int32(bits&0x7E000000) >> 20) | (int32(bits&0x00000F00) >> 7)
	inst.ImmJ = (int32(bits&0x80000000) >> 11) | int32(bits&0x000FF000) | (int32(bits&0x00100000) >> 9) | (int32(bits&0x7FE00000) >> 20)
	inst.ImmU = int32(bits&0xFFFFF000) >> 12
	return
}

func LookupInstruction(inst *Instruction) InstructionType {
	panic("unimplemented")
}
