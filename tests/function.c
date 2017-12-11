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

// Go keywords in C function
int chan()       {return 42;}
int defer()      {return 42;}
int fallthrough(){return 42;}
int func()       {return 42;}
int go()         {return 42;}
int import()     {return 42;}
int interface()  {return 42;}
int map()        {return 42;}
int package()    {return 42;}
int range()      {return 42;}
int select()     {return 42;}
int type()       {return 42;}
int var()        {return 42;}
int _()          {return 42;}

int main()
{
    plan(25);

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
	
	diag("Not allowable function name for Go")
	is_eq( chan()       , 42);
	is_eq( defer()      , 42);
	is_eq( fallthrough(), 42);
	is_eq( func()       , 42);
	is_eq( go()         , 42);
	is_eq( import()     , 42);
	is_eq( interface()  , 42);
	is_eq( map()        , 42);
	is_eq( package()    , 42);
	is_eq( range()      , 42);
	is_eq( select()     , 42);
	is_eq( type()       , 42);
	is_eq( var()        , 42);
	is_eq( _()          , 42);

    done_testing();
}

void my_function()
{
    pass("%s", "Welcome to my function. Feel at home.");
}
