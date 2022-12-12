.text

main:
	li x6, 1
	li x7, 2
	jal x5, sum
	auipc x8, 0
	beq x8, x5, end

sum:
	add x6, x6, x7
	jalr x0, x5, 0

end:
	add x6, x6, x7