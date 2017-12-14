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

double (*f2)(int, float, double);

double add2 (int a, float b, double c) {
    return c+a+b;
}

double mul2 (int a, float b, double c) {
    return a*b*c;
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
    plan(30);

    pass("%s", "Main function.");

    my_function();

    pass("%s", "Back in function main.");

	diag("pointer on function. Part 1")
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

	diag("pointer on function. Part 2")
	f = add;
	int i = f(3,4);
	is_eq(i,7);
	f = mul;
	int j = f(3,4);
	is_eq(j,12);
	if( f(2,3) == 6 ){
		pass("Ok")
	}

	diag("pointer on function. Part 3")
	{
		f2 = add2;
		double temp_data;
		temp_data = f2(4,2,3);
		is_eq(temp_data,9);
		f2 = mul2;
		is_eq(f2(4,2,3),24);
		double ttt = f2(1,1,1);
		is_eq(ttt,1);
		if(add2(2,3,1) == 6.0){
			pass("Ok")
		}
	}
	
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
