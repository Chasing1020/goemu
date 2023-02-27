CROSS_COMPILE = riscv64-unknown-elf-
CFLAGS = -nostdlib -fno-builtin -march=rv64g -mabi=lp64 -g -Wall -Ttext=0x80000000
GDB = $(CROSS_COMPILE)gdb
CC  = $(CROSS_COMPILE)gcc
OBJCOPY = $(CROSS_COMPILE)objcopy
OBJDUMP = $(CROSS_COMPILE)objdump


HEXDUMP = hexdump

ELF_FILES = $(wildcard *.elf)
BIN_FILES = $(wildcard *.bin)

%.bin: %.elf
	$(OBJCOPY) -O binary $^ $@

%.elf: %.s
	${CC} $(CFLAGS) $^ -o $@

hexdump: $(ELF_FILES) $(BIN_FILES)
	$(HEXDUMP) -C $^

objdump: $(ELF_FILES) $(BIN_FILES)
	@$(OBJDUMP) -S -M no-aliases $^

clean:
	rm $(BIN_FILES) $(ELF_FILES)

.PHONY: hexdump objdump clean
