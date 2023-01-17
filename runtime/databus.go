package runtime

import (
	"fmt"
	"goemu/config"
	"goemu/hw/uart"
)

type Bus struct {
	Mem  *Memory
	Uart *uart.Uart
}

func (b *Bus) Load(addr, bytes uint64) (uint64, error) {
	switch {
	case addr >= config.KernelBase && addr <= config.KernelEnd:
		return b.Mem.Load(addr, bytes)
	case addr >= uart.Base && addr < uart.End:
		return b.Uart.Load(addr, bytes)
	default:
		return 0, fmt.Errorf("invalid memory address: %x", addr)
	}
}

func (b *Bus) Store(addr, bytes, data uint64) error {
	switch {
	case addr >= config.KernelBase && addr <= config.KernelEnd:
		return b.Mem.Store(addr, bytes, data)
	case addr >= uart.Base && addr < uart.End:
		return b.Uart.Store(addr, bytes, data)
	default:
		return fmt.Errorf("invalid memory address: %x", addr)
	}
}
