package runtime

type Level uint8

const (
	User       Level = 0b00
	Supervisor Level = 0b01
	Machine    Level = 0b11
)

// Machine-level CSRs
const (
	// Machine information registers
	mvendorid  = 0xF11 // Vendor ID
	marchid    = 0xF12 // Architecture ID
	mimpid     = 0xF13 // Implementation ID
	mhartid    = 0xF14 // Hardware thread ID
	mconfigptr = 0xF15 // Pointer to configuration data structure

	// Machine Trap Setup
	mstatus    = 0x300 // Machine status register
	misa       = 0x301 // ISA and extensions
	medeleg    = 0x302 // Machine exception delegation register
	mideleg    = 0x303 // Machine interrupt delegation register
	mie        = 0x304 // Machine interrupt-enable register
	mtvec      = 0x305 // Machine trap-handler base address
	mcounteren = 0x306 // Machine counter enable
	mstatush   = 0x310 // Additional machine status register, RV32 only

	// Machine Trap Handling
	mscratch = 0x340 // Scratch register for machine trap handlers
	mepc     = 0x341 // Machine exception program counter
	mcause   = 0x342 // Machine trap cause
	mtval    = 0x343 // Machine bad address or instruction
	mip      = 0x344 // Machine interrupt pending
	mtinst   = 0x34A // Machine trap instruction (transformed)
	mtval2   = 0x34B // Machine bad guest physical address

)

// Supervisor-level CSRs
const (
	// Supervisor Trap Setup
	sstatus    = 0x100 // Supervisor status
	sedeleg    = 0x102 // Supervisor exception delegation
	sideleg    = 0x103 // Supervisor interrupt delegation
	sie        = 0x104 // Supervisor interrupt enable
	stvec      = 0x105 // Supervisor trap-handler base address
	scounteren = 0x106 // Supervisor counter enable

	// Supervisor Trap Handling
	sscratch = 0x140 // Supervisor scratch
	sepc     = 0x141 // Supervisor exception program counter
	scause   = 0x142 // Supervisor trap cause
	stval    = 0x143 // Supervisor trap value
	sip      = 0x144 // Supervisor interrupt-pending

	// Supervisor Protection and Translation
	satp = 0x180 // Supervisor address translation and protection
)
