package uart

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// UART control registers.
// ref to http://byterunner.com/16550.html
const (
	Base = 0x10000000
	Size = 0x100
	End  = Base + Size - 1

	Irq = 10 // UartInterrupt Request

	Rhr           = 0b000 // receive holding register (for input bytes)
	Thr           = 0b000 // transmit holding register (for output bytes)
	Ier           = 0b001 // interrupt enable register
	IerRxEnable   = 1 << 0
	IerTxEnable   = 1 << 1
	Fcr           = 0b010 // FIFO control register
	FcrFifoEnable = 1 << 0
	FcrFifoClear  = 3 << 1 // clear the content of the two FIFOs
	Isr           = 0b010  // interrupt status register
	Lcr           = 0b011  // line control register
	LcrEightBits  = 3 << 0
	LcrBaudLatch  = 1 << 7 // special mode to set baud rate
	Mcr           = 0b100
	Lsr           = 0b101
	LsrRxReady    = 1 << 0 // input is waiting to be read from RHR
	LsrTxIdle     = 1 << 5 // THR can accept another character to send
)

// 268435456

const BufferMaxSize = 32

type Uart struct {
	Regs [Size]uint8
	buf  strings.Builder

	in *bufio.Reader

	loadEnable chan struct{}
	mu         sync.Mutex
}

func NewUart() *Uart {
	u := new(Uart)
	u.Regs[Lsr] |= LsrTxIdle
	u.in = bufio.NewReader(os.Stdin)

	go u.InputHandler()

	return u
}

func (u *Uart) InputHandler() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	for {
		b, err := u.in.ReadByte()
		if err != nil {
			panic(err)
		}

		<-u.loadEnable // fixme: unknown output
		u.mu.Lock()
		u.Regs[Rhr] = b
		u.Regs[Lsr] |= LsrRxReady
		u.mu.Unlock()
	}
}

func (u *Uart) Check(bytes uint64) error {
	if bytes != 1 {
		return fmt.Errorf("invalid data bytes: %d", bytes)
	}
	return nil
}

func (u *Uart) Load(addr, bytes uint64) (uint64, error) {
	if err := u.Check(bytes); err != nil {
		return 0, err
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	switch addr - Base {
	case Rhr:
		u.Regs[Lsr] &= ^uint8(LsrRxReady)
		if u.Regs[Lsr]&LsrRxReady == 1 {
			u.loadEnable <- struct{}{}
		}
		return uint64(u.Regs[Rhr]), nil
	default:
		return uint64(u.Regs[addr-Base]), nil
	}
}

func (u *Uart) Store(addr, bytes, data uint64) error {
	if err := u.Check(bytes); err != nil {
		return err
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	r := rune(data)
	switch addr - Base {
	case Thr:
		if _, err := u.buf.WriteRune(r); err != nil {
			return err
		}
		if r == '\n' || u.buf.Len() >= BufferMaxSize {
			u.flushBuffer()
		}
		return nil
	default:
		u.Regs[addr-Base] = uint8(data)
		return nil
	}
}

func (u *Uart) flushBuffer() {
	fmt.Print(u.buf.String())
	u.buf.Reset()
}
