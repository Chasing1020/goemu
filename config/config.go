package config

const (
	MemSize = 128 * 1024 * 1024 // 128MiB

	KernelBase = 0x80000000
	KernelEnd  = KernelBase + MemSize - 1
)
