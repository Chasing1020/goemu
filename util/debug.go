package util

import "fmt"

func ParseInst(inst uint64) {
	opcode := uint8(inst & 0x0000007F)
	rs1 := uint8((inst & 0x000F8000) >> 15)
	rs2 := uint8((inst & 0x01F00000) >> 20)
	rd := uint8((inst & 0x00000F80) >> 7)
	funct3 := uint8((inst & 0x00007000) >> 12)
	funct7 := uint8((inst & 0xFE000000) >> 25)
	csr := (inst & 0xFFF00000) >> 20

	immI := uint64(int32(inst&0xFFF00000) >> 20)
	immS := (inst&0xFE000000)>>20 | (inst & 0x00000F80 >> 7)
	immB := uint64(int32(inst&0x80000000)>>19) | (inst & 0x80 << 4) | (inst >> 20 & 0x7E0) | (inst >> 7 & 0x1E)
	immJ := uint64((int32(uint64(inst)&0x80000000))>>11) | (uint64(inst) & 0xFF000) | ((inst >> 9) & 0x800) | ((inst >> 20) & 0x7FE)
	immU := inst & 0xFFFFF000
	fmt.Printf("%s: %08x(%032b)\n", "inst  ", inst, inst)
	fmt.Printf("%s: %02x(%07b)\n", "opcode", opcode, opcode)
	fmt.Printf("%s: %02x(%05b)\n", "rs1   ", rs1, rs1)
	fmt.Printf("%s: %02x(%05b)\n", "rs2   ", rs2, rs2)
	fmt.Printf("%s: %02x(%05b)\n", "rd    ", rd, rd)
	fmt.Printf("%s: %02x(%03b)\n", "funct3", funct3, funct3)
	fmt.Printf("%s: %02x(%07b)\n", "funct7", funct7, funct7)
	fmt.Printf("%s: %08x(%032b)\n", "csr ", csr, csr)
	fmt.Printf("%s: %08x(%032b)\n", "immI", immI, immI)
	fmt.Printf("%s: %08x(%032b)\n", "immS", immS, immS)
	fmt.Printf("%s: %08x(%032b)\n", "immB", immB, immB)
	fmt.Printf("%s: %08x(%032b)\n", "immJ", immJ, immJ)
	fmt.Printf("%s: %08x(%032b)\n", "immU", immU, immU)
}
