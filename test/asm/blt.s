.text
.global	_start

target:
    li x8, 1

_start:
    auipc x7, 0
    li x5, 10
    li x6, -10
    blt  x5, x6, target
