// Tests for unions.

#include <stdio.h>
#include "tests.h"


union simple_union
{
	int a;
	int b;
};

void union_simple()
{
	union simple_union s;
	s.a = 43;
	is_eq( s.a , 43 );
	is_eq( s.b , 43 );
}

int main()
{
    plan(2);

	union_simple();

    done_testing();
}
