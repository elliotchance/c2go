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

void test_castbool()
{
    char i1 = (1 == 1);
    short i2 = (1 == 1);
    int i3 = (1 == 1);
    long i4 = (1 == 0);
    long long i5 = (1 == 0);

    is_eq((i1==1) && (i2==1) && (i3==1) && (i4==0) && (i5==0), 1);
}

int main()
{
    plan(3);

    START_TEST(cast)
    START_TEST(castbool)

    done_testing();
}
