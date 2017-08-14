#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include "tests.h"

void test_malloc1()
{
    diag("malloc1");

    int i = 16, n;
    char *buffer;

    buffer = (char *)malloc(i + 1);
    is_not_null(buffer) or_return();

    for (n = 0; n < i; n++)
        buffer[n] = i % 26 + 'a';
    buffer[i] = '\0';

    is_streq(buffer, "qqqqqqqqqqqqqqqq");
    free(buffer);
}

void test_malloc2()
{
    diag("malloc2");

    int *p;
    p = (int *)malloc(sizeof(int));
    is_not_null(p) or_return();

    *p = 5;

    is_eq(*p, 5);
}

// Mix around all the types to make sure it is still actually allocating the
// correct size.
void test_malloc3()
{
    diag("malloc3");

    is_eq(sizeof(int), 4);
    is_eq(sizeof(double), 8);

    // 10 ints, should be 5 doubles. Also use a bad cast to make sure that it
    // doesn't interfere with the types.
    double *d;
    d = (char *)malloc(sizeof(int) * 10);
    is_not_null(d) or_return();

    // We can't test how much memory has been allocated by Go, but we can test
    // that there are 5 items.
    *d = 123;
    d[4] = 456;

    is_eq(*d, 123);
    is_eq(d[4], 456);
}

// calloc() works exactly the same as malloc() however the memory is zeroed out.
// In Go all allocated memory is zeroed out so they actually are the same thing.
void test_calloc()
{
    diag("calloc");

    is_eq(sizeof(int), 4);
    is_eq(sizeof(double), 8);

    // 10 ints, should be 5 doubles. Also use a bad cast to make sure that it
    // doesn't interfere with the types.
    double *d;
    d = (char *)calloc(sizeof(int), 10);
    is_not_null(d) or_return();

    // We can't test how much memory has been allocated by Go, but we can test
    // that there are 5 items.
    *d = 123;
    d[4] = 456;

    is_eq(*d, 123);
    is_eq(d[4], 456);
}

int main()
{
    plan(17);

    diag("abs")
    is_eq(abs(-5), 5);
    is_eq(abs(7), 7);
    is_eq(abs(0), 0);

    test_malloc1();
    test_malloc2();
    test_malloc3();
    test_calloc();

    done_testing();
}
