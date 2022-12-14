package config

const (
	MemSize = 1024 * 1024 * 1024

	KernelBase = 0x80000000
	KernelEnd  = KernelBase + MemSize - 1
)

var (
	AbiMap = [32]string{
		"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2",
		"s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
		"a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
		"s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
	}
)
