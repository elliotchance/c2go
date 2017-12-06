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
    plan(11);

    START_TEST(cast)
    START_TEST(castbool)

	double *d = (double *) 0;
	is_true(d == NULL);
	int    *i = (int    *) 0;
	is_true(i == NULL);
	float  *f = (float  *) 0;
	is_true(f == NULL);
	char   *c = (char   *) 0;
	is_true(c == NULL);

	double *d2 = 0;
	is_true(d2 == NULL);
	int    *i2 = 0;
	is_true(i2 == NULL);
	float  *f2 = 0;
	is_true(f2 == NULL);
	char   *c2 = 0;
	is_true(c2 == NULL);

    done_testing();
}
