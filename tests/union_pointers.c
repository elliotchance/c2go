// Tests for unions.

#include <stdio.h>
#include "tests.h"

union UPNT{
	int * a;
	int * b;
	int * c;
};

void union_pointers()
{
	union UPNT u;
	int v = 32;
	u.a = &v;
	is_eq(*u.a,32);
	is_eq(*u.b,32);
	is_eq(*u.c,32);
}

union UPNTF{
	int (*f1)(int);
	int (*f2)(int);
};

int union_function(int a)
{
	return a+1;
}

void union_func_pointers()
{
	union UPNTF u;
	u.f1 = union_function;
	is_eq(u.f1(21), 22);
	is_eq(u.f2(21), 22);
}


int main()
{
    plan(5);
	
	union_pointers();
	union_func_pointers();
    
	done_testing();
}
