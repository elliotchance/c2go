// Array examples

#include <stdio.h>
#include "tests.h"

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_intarr()
{
    int a[3];
    a[0] = 5;
    a[1] = 9;
    a[2] = -13;

    is_eq(a[0], 5);
    is_eq(a[1], 9);
    is_eq(a[2], -13);
}

void test_doublearr()
{
    double a[2];
    a[0] = 1.2;
    a[1] = 7; // different type

    is_eq(a[0], 1.2);
    is_eq(a[1], 7.0);
}

void test_intarr_init()
{
    int a[] = {10, 20, 30};
    is_eq(a[0], 10);
    is_eq(a[1], 20);
    is_eq(a[2], 30);
}

void test_floatarr_init()
{
    float a[] = {2.2, 3.3, 4.4};
    is_eq(a[0], 2.2);
    is_eq(a[1], 3.3);
    is_eq(a[2], 4.4);
}

void test_chararr_init()
{
    char a[] = {97, 98, 99};
    is_eq(a[0], 'a');
    is_eq(a[1], 'b');
    is_eq(a[2], 'c');
}

void test_chararr_init2()
{
    char a[] = {'a', 'b', 'c'};
    is_eq(a[0], 'a');
    is_eq(a[1], 'b');
    is_eq(a[2], 'c');
}

void test_exprarr()
{
    int a[] = {2 ^ 1, 3 & 1, 4 | 1, (5 + 1)/2};
    is_eq(a[0], 3);
    is_eq(a[1], 1);
    is_eq(a[2], 5);
    is_eq(a[3], 3);
}

int main()
{
    plan(21);

    START_TEST(intarr);
    START_TEST(doublearr);
    START_TEST(intarr_init);
    START_TEST(floatarr_init);
    START_TEST(chararr_init);
    START_TEST(chararr_init2);
    START_TEST(exprarr);

    done_testing();
}
