// This file contains tests for the sizeof() function and operator.

#include <stdio.h>
#include "tests.h"

#define check_sizes(type, size)         \
    is_eq(sizeof(type), size);          \
    is_eq(sizeof(unsigned type), size); \
    is_eq(sizeof(signed type), size);   \
    is_eq(sizeof(const type), size);    \
    is_eq(sizeof(volatile type), size);

#define FLOAT(type, size) \
    is_eq(sizeof(type), size);

#define OTHER(type, size) \
    is_eq(sizeof(type), size);

// We print the variable so that the compiler doesn't complain that the variable
// is unused.
#define VARIABLE(v, p) \
    printf("%s = (%d) %d bytes\n", #v, p, sizeof(v));

struct MyStruct
{
    double a;
    char b;
    char c;
};

short a;
int b;

int main()
{
    plan(32);

    diag("Integer types");
    check_sizes(char, 1);
    check_sizes(short, 2);
    check_sizes(int, 4);
    check_sizes(long, 8);

    diag("Floating-point types");
    is_eq(sizeof(float), 4);
    is_eq(sizeof(double), 8);
    is_eq(sizeof(long double), 16);

    diag("Other types");
    is_eq(sizeof(void), 1);

    diag("Pointers");
    is_eq(sizeof(char *), 8);
    is_eq(sizeof(char *), 8);
    is_eq(sizeof(short **), 8);

    diag("Variables");
    a = 123;
    b = 456;
    struct MyStruct s1;
    s1.b = 0;

    is_eq(sizeof(a), 2);
    is_eq(sizeof(b), 4);
    is_eq(sizeof(s1), 16);

    diag("Structures");
    is_eq(sizeof(struct MyStruct), 16);

    diag("Function pointers");
    is_eq(sizeof(main), 1);

    diag("TODO: Unions");

    done_testing();
}
