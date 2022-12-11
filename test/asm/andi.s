.text
.global	_start

_start:
	li x6, 0x91  # x6 = 0b1001_0001
	andi x5, x6, 0x01	# x5 = x6 & 0x01
