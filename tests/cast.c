#include <stdio.h>
#include "tests.h"

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_cast()
{
    int c[] = {(int)'a', (int)'b'};
    is_eq(c[0], 97);

    double d = (double) 1;
    is_eq(d, 1.0);
}

int main()
{
    plan(2);

    START_TEST(cast)

    done_testing();
}
