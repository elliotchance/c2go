// This file contains tests for the sizeof() function and operator.

#include <stdio.h>

#define T(type) printf("%s = %d bytes\n", #type, sizeof(type))

int main(int argc, char *argv[])
{
    T(int);

    return 0;
}
