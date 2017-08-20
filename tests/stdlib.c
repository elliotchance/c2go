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
    plan(108);

    diag("abs")
    is_eq(abs(-5), 5);
    is_eq(abs(7), 7);
    is_eq(abs(0), 0);

    diag("atof")
    is_eq(atof("123"), 123);
    is_eq(atof("1.23"), 1.23);
    is_eq(atof(""), 0);
    is_eq(atof("1.2e6"), 1.2e6);
    is_eq(atof(" \n123"), 123);
    is_eq(atof("\t123foo"), 123);
    is_eq(atof("+1.23"), 1.23);
    is_eq(atof("-1.23"), -1.23);
    is_eq(atof("1.2E-6"), 1.2e-6);
    is_eq(atof("1a2b"), 1);
    is_eq(atof("1a.2b"), 1);
    is_eq(atof("a1.2b"), 0);
    is_eq(atof("1.2Ee-6"), 1.2);
    is_eq(atof("-1..23"), -1);
    is_eq(atof("-1.2.3"), -1.2);
    is_eq(atof("foo"), 0);
    is_eq(atof("+1.2+3"), 1.2);
    is_eq(atof("-1.-23"), -1);
    is_eq(atof("-.23"), -0.23);
    is_eq(atof(".4"), 0.4);
    is_eq(atof("0xabc"), 2748);
    is_eq(atof("0x1b9"), 441);
    is_eq(atof("0x"), 0);
    is_eq(atof("0X1f9"), 505);
    is_eq(atof("-0X1f9"), -505);
    is_eq(atof("+0x1f9"), 505);
    is_eq(atof("0X"), 0);
    is_eq(atof("0xfaz"), 250);
    is_eq(atof("0Xzaf"), 0);
    is_eq(atof("0xabcp2"), 10922);
    is_eq(atof("0xabcP3"), 10922);
    is_eq(atof("0xabcP2z"), 10922);
    is_eq(atof("0xabcp-2"), 687);
    is_eq(atof("0xabcp+2"), 10922);
    is_inf(atof("inf"), 1);
    is_inf(atof("INF"), 1);
    is_inf(atof("Inf"), 1);
    is_inf(atof("-Inf"), -1);
    is_inf(atof("+INF"), 1);
    is_inf(atof("infinity"), 1);
    is_inf(atof("INFINITY"), 1);
    is_inf(atof("Infinity"), 1);
    is_inf(atof("+INFINITY"), 1);
    is_inf(atof("-InfINITY"), -1);
    is_nan(atof("nan"));
    is_nan(atof("NaN"));
    is_nan(atof("+NaN"));
    is_nan(atof("NAN"));
    is_nan(atof("-NAN"));
    is_nan(atof("nanabc123"));
    is_nan(atof("NANz123"));
    is_nan(atof("NaN123z"));
    is_nan(atof("-NANz123"));
    // This causes a segfault in C:
    // is_eq(atof(NULL), 0);

    diag("atoi")
    is_eq(atoi("123"), 123)
    is_eq(atoi("+456"), 456)
    is_eq(atoi("-52"), -52)
    is_eq(atoi("foo"), 0)
    is_eq(atoi(" \t123"), 123)
    is_eq(atoi(""), 0)
    is_eq(atoi(" \t"), 0)
    is_eq(atoi("123abc"), 123)

    diag("atol")
    is_eq(atol("123"), 123)
    is_eq(atol("+456"), 456)
    is_eq(atol("-52"), -52)
    is_eq(atol("foo"), 0)
    is_eq(atol(" \t123"), 123)
    is_eq(atol(""), 0)
    is_eq(atol(" \t"), 0)
    is_eq(atol("123abc"), 123)

    diag("atoll")
    is_eq(atoll("123"), 123)
    is_eq(atoll("+456"), 456)
    is_eq(atoll("-52"), -52)
    is_eq(atoll("foo"), 0)
    is_eq(atoll(" \t123"), 123)
    is_eq(atoll(""), 0)
    is_eq(atoll(" \t"), 0)
    is_eq(atoll("123abc"), 123)

    diag("div")
    div_t result = div(17, 5);
    is_eq(result.quot, 3)
    is_eq(result.rem, 2)
    result = div(-17, 5);
    is_eq(result.quot, -3)
    is_eq(result.rem, -2)
    result = div(17, -5);
    is_eq(result.quot, -3)
    is_eq(result.rem, 2)
    result = div(-17, -5);
    is_eq(result.quot, 3)
    is_eq(result.rem, -2)

    // exit() is handled in tests/exit.c

    // free() is handled with the malloc and calloc tests.

    diag("getenv")
    is_not_null(getenv("PATH"));
    is_not_null(getenv("HOME"));
    is_null(getenv("FOOBAR"));

    diag("labs")
    is_eq(labs(-5), 5);
    is_eq(labs(7), 7);
    is_eq(labs(0), 0);

    test_malloc1();
    test_malloc2();
    test_malloc3();
    test_calloc();

    done_testing();
}
