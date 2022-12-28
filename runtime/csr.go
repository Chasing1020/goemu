package runtime

import "fmt"

type Level uint8

const (
	User       Level = 0b00
	Supervisor Level = 0b01
	Machine    Level = 0b11
)

// Machine Level CSRs
const (
	// Machine information registers
	Mvendorid  = 0xF11 // Vendor ID
	Marchid    = 0xF12 // Architecture ID
	Mimpid     = 0xF13 // Implementation ID
	Mhartid    = 0xF14 // Hardware thread ID
	Mconfigptr = 0xF15 // Pointer to configuration data structure

	// Machine Trap Setup
	Mstatus    = 0x300 // Machine status register
	Misa       = 0x301 // ISA and extensions
	Medeleg    = 0x302 // Machine exception delegation register
	Mideleg    = 0x303 // Machine interrupt delegation register
	Mie        = 0x304 // Machine interrupt-enable register
	Mtvec      = 0x305 // Machine trap-handler base address
	Mcounteren = 0x306 // Machine counter enable
	Mstatush   = 0x310 // Additional machine status register, RV32 only

	// Machine Trap Handling
	Mscratch = 0x340 // Scratch register for machine trap handlers
	Mepc     = 0x341 // Machine exception program counter
	Mcause   = 0x342 // Machine trap cause
	Mtval    = 0x343 // Machine bad address or instruction
	Mip      = 0x344 // Machine interrupt pending
	Mtinst   = 0x34A // Machine trap instruction (transformed)
	Mtval2   = 0x34B // Machine bad guest physical address
)

// Supervisor Level CSRs
const (
	// Supervisor Trap Setup
	Sstatus    = 0x100 // Supervisor status
	Sedeleg    = 0x102 // Supervisor exception delegation
	Sideleg    = 0x103 // Supervisor interrupt delegation
	Sie        = 0x104 // Supervisor interrupt enable
	Stvec      = 0x105 // Supervisor trap-handler base address
	Scounteren = 0x106 // Supervisor counter enable

	// Supervisor Trap Handling
	Sscratch = 0x140 // Supervisor scratch
	Sepc     = 0x141 // Supervisor exception program counter
	Scause   = 0x142 // Supervisor trap cause
	Stval    = 0x143 // Supervisor trap value
	Sip      = 0x144 // Supervisor interrupt-pending

	// Supervisor Protection and Translation
	Satp = 0x180 // Supervisor address translation and protection
)

// Mstatus and Sstatus field mask
const (
	SieMask            = 1 << 1
	MieMask            = 1 << 3
	SpieMask           = 1 << 5
	UbeMask            = 1 << 6
	MpieMask           = 1 << 7
	SppMask            = 1 << 8
	VsMask             = 0b11 << 9
	MppMask            = 0b11 << 11
	FsMask             = 0b11 << 13
	XsMask             = 0b11 << 15
	MprvMask           = 1 << 17
	SumMask            = 1 << 18
	MxrMask            = 1 << 19
	TvmMask            = 1 << 20
	TwMask             = 1 << 21
	TsrMask            = 1 << 22
	UxlMask            = 0b11 << 32
	SxlMask            = 0b11 << 34
	SbeMask            = 1 << 36
	MbeMask            = 1 << 37
	SdMask             = 1 << 63
	SstatusMask uint64 = SieMask | SpieMask | UbeMask | SppMask | FsMask | XsMask | SumMask | MxrMask | UxlMask | SdMask
)

const CsrNum = 0xFFF + 1

type CSR [CsrNum]uint64

func (c *CSR) Load(addr uint64) (uint64, error) {
	if addr > CsrNum || addr < 0 {
		return 0, fmt.Errorf("invalid csr address: %x", addr)
	}

	switch addr {
	case Sie:
		return (*c)[Mie] & (*c)[Mideleg], nil
	case Sip:
		return (*c)[Mip] & (*c)[Mideleg], nil
	case Sstatus:
		return (*c)[Mstatus] & SstatusMask, nil
	default:
		return (*c)[addr], nil
	}
}

func (c *CSR) Store(addr, data uint64) error {
	if addr > CsrNum || addr < 0 {
		return fmt.Errorf("invalid csr address: %x", addr)
	}

	switch addr {
	case Sie:
		(*c)[Mie] = ((*c)[Mie] & ^(*c)[Mideleg]) | (data & (*c)[Mideleg])
	case Sip:
		(*c)[Mip] = ((*c)[Mip] & ^(*c)[Mideleg]) | (data & (*c)[Mideleg])
	case Sstatus:
		(*c)[Mstatus] = ((*c)[Mstatus] & ^SstatusMask) | (data & SstatusMask)
	default:
		(*c)[addr] = data
	}
	return nil
}
