main:
	li x6, 0x12313	# x6 = 0x12313

	not x5, x6		    # x5 = ~x6
	xori x7, x6, -1		# the same as not