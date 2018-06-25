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

struct MyStruct2
{
    double a;
    char b;
    char c;
    char d[10];
};

struct MyStruct3
{
    double a;
    char b;
    char c;
    char d[20];
};

struct MyStruct4
{
    double a;
    char b;
    char c;
    char d[30];
};

union MyUnion
{
    double a;
    char b;
    int c;
};

short a;
int b;

int main()
{
    plan(42);

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
    struct MyStruct2 s2;
    s2.b = 0;
    struct MyStruct3 s3;
    s3.b = 0;
    struct MyStruct4 s4;
    s4.b = 0;
    union MyUnion u1;
    u1.b = 0;

    is_eq(sizeof(a), 2);
    is_eq(sizeof(b), 4);
    is_eq(sizeof(s1), 16);
    is_eq(sizeof(s2), 24);
    is_eq(sizeof(s3), 32);
    is_eq(sizeof(s4), 40);
    is_eq(sizeof(u1), 8);

    diag("Structures");
    is_eq(sizeof(struct MyStruct), 16);

    diag("Unions");
    is_eq(sizeof(union MyUnion), 8);

    diag("Function pointers");
    is_eq(sizeof(main), 1);

    diag("Arrays");
    char c[3] = {'a', 'b', 'c'};
    c[0] = 'a';
    is_eq(sizeof(c), 3);

    int *d[3];
    d[0] = &b;
    is_eq(sizeof(d), 24);

    int **e[4];
    e[0] = d;
    is_eq(sizeof(e), 32);

    const char * const f[] = {"a", "b", "c", "d", "e", "f"};
    is_eq(sizeof(f), 48);
    is_streq(f[1], "b");

    done_testing();
}
