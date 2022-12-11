package runtime

// InstructionFields includes fields for all possible instruction types
// can be a convenient way to represent RISC-V instructions
// not all fields will be valid for all instruction types
type InstructionFields struct {
	Bits uint32 // for debugging

	Opcode uint8
	Rs1    uint8
	Rs2    uint8
	Rd     uint8
	Funct3 uint8
	Funct7 uint8
	Csr    uint32

	ImmI uint64
	ImmS uint64
	ImmB uint64
	ImmJ uint64
	ImmU uint64
}

// ParseInstruction extracts the instruction's opcode, register fields, and immediate values by masking and shifting the input Bits.
func ParseInstruction(bits uint32) (inst InstructionFields) {
	inst.Bits = bits

	inst.Opcode = uint8(bits & 0x0000007F)
	inst.Rs1 = uint8((bits & 0x000F8000) >> 15)
	inst.Rs2 = uint8((bits & 0x01F00000) >> 20)
	inst.Rd = uint8((bits & 0x00000F80) >> 7)
	inst.Funct3 = uint8((bits & 0x00007000) >> 12)
	inst.Funct7 = uint8((bits & 0xFE000000) >> 25)
	inst.Csr = (bits & 0xFFF00000) >> 20

	inst.ImmI = uint64(bits&0xFFF00000) >> 20
	inst.ImmS = (uint64(bits&0xFE000000) >> 20) | (uint64(bits&0x00000F80) >> 7)
	inst.ImmB = (uint64(bits&0x80000000) >> 19) | (uint64(bits&0x00000080) << 4) | (uint64(bits&0x7E000000) >> 20) | (uint64(bits&0x00000F00) >> 7)
	inst.ImmJ = (uint64(bits&0x80000000) >> 11) | uint64(bits&0x000FF000) | (uint64(bits&0x00100000) >> 9) | (uint64(bits&0x7FE00000) >> 20)
	inst.ImmU = uint64(bits&0xFFFFF000) >> 12
	return
}
