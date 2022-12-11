.data
.word 2, 4, 6, 8
n: .word 5

.text
main:
    # add
    add t0, x0, x0
    addi t1, x0, 1
    la t3, n
    lw t3, 0(t3)
fib:
    beq t3, x0, finish
    add t2, t1, t0
    mv t0, t1
    mv t1, t2
    addi t3, t3, -1
    j fib
finish:
    addi a0, x0, 1 # stdout
    addi a1, t0, 0 # a1 should be 8
