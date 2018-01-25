// Tests for unions.

#include <stdio.h>
#include "tests.h"

union int_union
{
	int a;
	int b;
};

union float_union
{
	float a;
	float b;
};

void union_simple()
{
	union int_union i;
	i.a = 43;
	is_eq( i.a , 43 );
	is_eq( i.b , 43 );

	union float_union f;
	f.a = 45;
	is_eq( f.a , 45 );
	is_eq( f.b , 45 );
}

union array_union
{
	float a[2];
	float b[2];
};

void union_array()
{
	union array_union arr;
	arr.a[0] = 12;
	arr.b[1] = 14;
	is_eq( arr.a[0] , 12);
	is_eq( arr.a[1] , 14);
	is_eq( arr.b[0] , 12);
	is_eq( arr.b[1] , 14);
}

int main()
{
    plan(8);

	union_simple();
	union_array ();

    done_testing();
}
