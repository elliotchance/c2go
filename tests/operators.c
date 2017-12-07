#include <stdio.h>
#include "tests.h"

// TODO: More comprehensive operator tests
// https://github.com/elliotchance/c2go/issues/143

int main()
{
	plan(49);

    int i = 10;
    signed char j = 1;
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

	diag("Shift with signed int");

    i = 4 << j;
        is_eq(i, 8);

    i = 8 >> j;
        is_eq(i, 4);

    i <<= j;
        is_eq(i, 8);

    i >>= j;
        is_eq(i, 4);

	diag("Operator equal for 1 variables");
	int x;
	x = 42;
		is_eq(x, 42);

	diag("Operator equal for 2 variables");
	int y;
	x = y = 1;
		is_eq(x, 1);
		is_eq(y, 1);
	
	diag("Operator comma in initialization");
	int x2 = 0, y2 = 1;
		is_eq(x2, 0);
		is_eq(y2, 1);

	diag("Operator equal for 3 variables");
	int a,b,c2;
	a = b = c2 = 3;
		is_eq(a, 3);
		is_eq(b, 3);
		is_eq(c2, 3);
	
	diag("Huge comma problem for Equal operator")
	int q,w,e;
	q = 7, w = q + 3, e = q + w;
		is_eq(q, 7);
		is_eq(w, 10);
		is_eq(e, 17);

	diag("Huge comma problem for Equal operator with Multiplication")
	float qF,wF,eF;
	qF = 7., wF = qF * 3., eF = qF * wF;
	float expectedQ = 7.;
	float expectedW = 7. * 3.;
	float expectedE = 7. * (7. * 3.);
		is_eq(qF, expectedQ); 
		is_eq(wF, expectedW);
		is_eq(eF, expectedE);

	diag("Statement expressions")
	int s1 = ({ 2; });
	is_eq(s1, 2);
	is_eq(({ int foo = s1 * 3; foo + 1; }), 7);

	diag("Not allowable var name for Go")
	int type = 42;
	is_eq(type,42);

	diag("Go keywords inside C code")
	{
		int chan = 42;
		is_eq(chan ,42);
	}
	{
		int defer = 42;
		is_eq(defer ,42);
	}
	{
		int fallthrough = 42;
		is_eq(fallthrough ,42);
	}
	{
		int func = 42;
		is_eq(func ,42);
	}
	{
		int go = 42;
		is_eq(go ,42);
	}
	{
		int import = 42;
		is_eq(import ,42);
	}
	{
		int interface = 42;
		is_eq(interface ,42);
	}
	{
		int map = 42;
		is_eq(map ,42);
	}
	{
		int package = 42;
		is_eq(package ,42);
	}
	{
		int range = 42;
		is_eq(range ,42);
	}
	{
		int select = 42;
		is_eq(select ,42);
	}
	{
		int type = 42;
		is_eq(type ,42);
	}
	{
		int var = 42;
		is_eq(var ,42);
	}
	
	done_testing();
}
