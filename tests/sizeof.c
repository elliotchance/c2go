// This file contains tests for the sizeof() function and operator.

#include <stdio.h>

#define INT(type) \
    printf("%s = %d bytes\n", #type, sizeof(type)); \
    printf("unsigned %s = %d bytes\n", #type, sizeof(unsigned type)); \
    printf("signed %s = %d bytes\n", #type, sizeof(signed type));

#define FLOAT(type) \
    printf("%s = %d bytes\n", #type, sizeof(type));

#define OTHER(type) \
    printf("%s = %d bytes\n", #type, sizeof(type));

// We print the variable so that the compiler doesn't complain that the variable
// is unused.
#define VARIABLE(v) \
    printf("%s = (%d) %d bytes\n", #v, v, sizeof(v));

int main(int argc, char *argv[])
{
    // Integer types.
    INT(char)
    INT(short)
    INT(int)
    INT(long)

    // Floating-point types.
    FLOAT(float)
    FLOAT(double)
    FLOAT(long double)

    // Other types.
    OTHER(void)

    // Types with qualifiers that do not effect the size.
    OTHER(const short)
    OTHER(volatile long double)

    // TODO: Pointers.
    OTHER(char*)
    OTHER(char *)
    OTHER(short**)

    // TODO: Variables.
    short a;
    int b;
    
    VARIABLE(a);
    VARIABLE(b);

    // TODO: Structures.
    // TODO: Function pointers.

    return 0;
}
