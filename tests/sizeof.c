// This file contains tests for the sizeof() function and operator.

#include <stdio.h>

#define INT(type) \
    printf("%s = %d bytes\n", #type, sizeof(type)); \
    printf("unsigned %s = %d bytes\n", #type, sizeof(unsigned type)); \
    printf("signed %s = %d bytes\n", #type, sizeof(signed type));

int main(int argc, char *argv[])
{
    // Basic types
    INT(char)
    INT(short)
    INT(int)
    INT(long)

    return 0;
}
