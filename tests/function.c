#include <stdio.h>
#include "tests.h"

void my_function();

int i = 40;

void function(){
	i += 2;
}

void function2(){
	i += 8;
}


int (*f)(int, int);

int add(int a, int b) {
        return a + b;
}

int mul(int a, int b) {
        return a * b;
}

int main()
{
    plan(11);

    pass("%s", "Main function.");

    my_function();

    pass("%s", "Back in function main.");

	// pointer on function
	void * a = NULL;
	is_null(a);
	a = function;
	is_not_null(a);
	void(*t)(void) = a;
	is_not_null(t);
	t();
	is_eq(i,42);
	t = function2;
	is_not_null(t);
	t();
	is_eq(i,50);

	// pointer on function
	f = add;
	int i = f(3,4);
	is_eq(i,7);
	f = mul;
	int j = f(3,4);
	is_eq(j,12);

    done_testing();
}

void my_function()
{
    pass("%s", "Welcome to my function. Feel at home.");
}
