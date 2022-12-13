.text

main:
    li x5, 0x12345678
    addi sp, sp, -16
    sd   x5, 8(sp)
    lb   x6, 8(sp)
    lh   x7, 8(sp)