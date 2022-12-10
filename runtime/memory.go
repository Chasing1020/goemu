package runtime

import (
	"fmt"
	"goemu/config"
)

type Memory []uint8

func (m *Memory) Check(addr, size uint64) error {
	if addr < config.KernelBase || addr > config.KernelEnd {
		return fmt.Errorf("invalid memory address: %x", addr)
	}
	switch size {
	case 8, 16, 32, 64:
		return nil
	default:
		return fmt.Errorf("invalid data size: %d", size)
	}
}

func (m *Memory) Load(addr, size uint64) (uint64, error) {
	if err := m.Check(addr, size); err != nil {
		return 0, err
	}
	index := addr - config.KernelBase
	data := uint64((*m)[index])
	for i := uint64(1); i < size/8; i++ {
		data |= uint64((*m)[index+i]) << (i * 8)
	}
	return data, nil
}

func (m *Memory) Store(addr, size uint64, data uint64) error {
	if err := m.Check(addr, size); err != nil {
		return err
	}
	index := addr - config.KernelBase
	for i := uint64(0); i < size/8; i++ {
		offset := 8 * i
		(*m)[index+i] = uint8((data >> offset) & 0xff)
	}
	return nil
}
