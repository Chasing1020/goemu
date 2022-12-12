.text


main:
	li x6, 0xff123ab	# int x6 = 0xffffffab
	la x5, _array		# array[0] = (char)x6
	sb x6, 0(x5)
    j end

_array:
	.byte 0x12
	.byte 0x34
    .byte 0x56
    .byte 0x78

end:
    nop