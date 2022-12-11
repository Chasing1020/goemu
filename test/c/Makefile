CROSS_COMPILE = riscv64-unknown-elf-
CFLAGS = -march=rv64g -ffreestanding -nostdlib -static -Wl,-Ttext=0 -O2
GDB = gdb-multiarch
CC  = ${CROSS_COMPILE}gcc
OBJCOPY = ${CROSS_COMPILE}objcopy
OBJDUMP = ${CROSS_COMPILE}objdump


HEXDUMP = hexdump

ELF_FILES = $(wildcard *.elf)
BIN_FILES = $(wildcard *.bin)

%.bin: %.elf
	${OBJCOPY} -O binary $^ $@

%.elf: %.c
	${CC} ${CFLAGS} $^ -o $@

.PHONY: hexdump
hexdump: ${ELF_FILES}
	${HEXDUMP} -C $^

.PHONY: objdump
objdump: ${ELF_FILES}
	@${OBJDUMP} -M no-aliases -d $^

.PHONY: clean
clean:
	rm ${BIN_FILES} ${ELF_FILES}