#include <stdio.h>
#include "tests.h"

// TODO: More comprehensive operator tests
// https://github.com/elliotchance/c2go/issues/143

int main()
{
	plan(15);

    int i = 10;
    float f = 3.14159f;
    double d = 0.0;
    char c = 'A';

    i %= 10;
	is_eq(i, 0);

    i += 10;
	is_eq(i, 10);

    i -= 2;
	is_eq(i, 8);

    i *= 2;
	is_eq(i, 16);
	
    i /= 4;
	is_eq(i, 4);

    i <<= 2;
	is_eq(i, 16);

    i >>= 2;
	is_eq(i, 4);

    i ^= 0xCFCF;
	is_eq(i, 53195);

    i |= 0xFFFF;
	is_eq(i, 65535);
	
    i &= 0x0000;
	is_eq(i, 0);
	
	diag("Other types");

    f += 1.0f;
	is_eq(f, 4.14159);

    d += 1.25f;
	is_eq(d, 1.25);
	
    i -= 255l;
	is_eq(i, -255);
	
    i += 'A';
	is_eq(i, -190);

    c += 11;
	is_eq(c, 76);

	done_testing();
}
