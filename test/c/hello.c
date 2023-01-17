void _start() {
    volatile char *p = (volatile char *)0x10000000;
    *p = 'z';
    *p = 'j';
    *p = 'c';
    *p = '\n';
}