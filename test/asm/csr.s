.text
main:
    li t2, 3
    csrrw zero, mstatus, 1
    csrrs zero, mtvec, 2
    csrrw zero, mepc, 3
    csrrc t2, mepc, 0
    csrrwi zero, sstatus, 4
    csrrsi zero, stvec, 5
    csrrwi zero, sepc, 5
    csrrci zero, sepc, 0
