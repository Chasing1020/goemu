package runtime

import "goemu/hw/uart"

type DataBus struct {
	Mem  Memory
	Uart uart.Uart
}
