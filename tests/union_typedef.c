// Tests for unions.

#include <stdio.h>
#include "tests.h"

typedef union myunion myunion;
typedef union myunion
{
	double PI;
	int B;
}MYUNION;

typedef union
{
	double PI;
	int B;
}MYUNION2;

void union_typedef()
{
	diag("Typedef union")
	union myunion m;
	double v = 3.14;
	m.PI = v;
	is_eq(m.PI,3.14);
	is_true(m.B != 0);
	is_eq(v, 3.14);
	v += 1.0;
	is_eq(v, 4.14);
	is_eq(m.PI,3.14);

	MYUNION mm;
	mm.PI = 3.14;
	is_eq(mm.PI,3.14);
	is_true(mm.B != 0);

	myunion mmm;
	mmm.PI = 3.14;
	is_eq(mmm.PI,3.14);
	is_true(mmm.B != 0);

	MYUNION2 mmmm;
	mmmm.PI = 3.14;
	is_eq(mmmm.PI,3.14);
	is_true(mmmm.B != 0);
}


int main()
{
    plan(11);
	
	union_typedef();
    
	done_testing();
}
