package runtime

import (
	"fmt"
	"goemu/config"
)

type Memory []uint8

func (m *Memory) Check(addr, bytes uint64) error {
	if addr < config.KernelBase || addr > config.KernelEnd {
		return fmt.Errorf("invalid memory address: %x", addr)
	}
	switch bytes {
	case 1, 2, 4, 8:
		return nil
	default:
		return fmt.Errorf("invalid data bytes: %d", bytes)
	}
}

func (m *Memory) Load(addr, bytes uint64) (uint64, error) {
	if err := m.Check(addr, bytes); err != nil {
		return 0, err
	}
	index := addr - config.KernelBase
	data := uint64((*m)[index])
	for i := uint64(1); i < bytes; i++ {
		data |= uint64((*m)[index+i]) << (i * 8)
	}
	return data, nil
}

func (m *Memory) Store(addr, bytes uint64, data uint64) error {
	if err := m.Check(addr, bytes); err != nil {
		return err
	}
	index := addr - config.KernelBase
	for i := uint64(0); i < bytes; i++ {
		offset := 8 * i
		(*m)[index+i] = uint8((data >> offset) & 0xFF)
	}
	return nil
}
