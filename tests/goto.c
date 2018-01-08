#include <stdio.h>
#include <time.h>

#include "tests.h"

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_goto1()
{
    int i = 0;

mylabel:
    i++;

    if (i > 5) {
        fail("Parameter i = %d, but expect 1", i);
        return;
    }

    if (i == 1) {
        goto mylabel;
    }

    is_eq(i, 2);
}

void test_goto2()
{
    int i = 0;
    int j = 0;

mylabel:
    i++, j++;

    if (j > 5) {
        fail("Parameter j = %d, but expect 1", j);
        return;
    }
    if (i == 1) {
        goto mylabel;
    }

    is_eq(i, 2);
    is_eq(j, 2);
}

void test_goto_stmt()
{
    int i = 0, j = 0;

mylabel:
    for (j = 0; j < 5; j++)
        i++;

    if (i < 15) {
        goto mylabel;
    }

    is_eq(i, 15);
}

int main()
{
    plan(4);

    START_TEST(goto1)
    START_TEST(goto2)
    START_TEST(goto_stmt)

    done_testing();
}
