.text
main:
	li x6, 1		# x6 = 1
	neg x5, x6		# x5 = -x6
	sub x7, x0, x6	# these two instructions assemble into the same thing!
