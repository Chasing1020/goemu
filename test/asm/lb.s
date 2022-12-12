.text

main:
    auipc x8, 0
    li x1, 0x1234
	sb x1, 0x100(x8)
	lb x7, 0x100(x8)
