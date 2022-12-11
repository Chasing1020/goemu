.text
.global	_start


_start:
	li x6, 1        # x6 = 1
	li x7, 2        # x7 = 2
	jal x5, sum		# call sum, return address is saved in x5
    auipc x8, 0     # x8 = PC

sum:
    add x6, x6, x7	# x6 = x6 + x7
    jalr x5, x0, 0 	# return
