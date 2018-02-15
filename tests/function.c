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

double action(double (*F)(int,float,double)){
	return F(2,3,4);
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
int init()       {return 42;}
int len()        {return 42;}
int copy()       {return 42;}
int fmt()        {return 42;}
int cap()        {return 42;}

void exit2(int t){
	(void)(t);
}
void empty_return(){
	int i = 0;
	(void)(i);
	exit2(-1);
}
int  empty_return_int(int a){
	if ( a > 0 ){ return 1;} 
	else        { exit2(-1);}
}
int empty_return_int2(int *a){
	if (*a > 0 ){ return 1;} 
	else        { exit2(-1);}
}
double empty_return_double(double a){
	if ( a > 0.0 ){ return 1.0;} 
	else          { exit2(-1);}
}
double empty_return_double2(double*a){
	if (*a > 0.0 ){ return 1.0;} 
	else          { exit2(-1);}
}
typedef struct RE RE;
struct RE {
	int re;
};
RE empty_return_struct(int a){
	if ( a > 0.0 ){ RE r; r.re = 1; return r;} 
	else          { exit2(-1);}
}
RE* empty_return_struct2(int a){
	if ( a > 0.0 ){ RE r; r.re = 1; return &r;} 
	else          { exit2(-1);}
}

typedef int (* operators)(int a,int b);
int call_a_func(operators call_this) {
    int output = call_this(5, 8);
    return output;
}

long tolower (int a, int b) { return (long)(a+b);}
long toupper (int a, int b) { return (long)(a+b);}

int main()
{
    plan(48);

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
	is_eq( init()       , 42);
	is_eq( len()        , 42);
	is_eq( copy()       , 42);
	is_eq( fmt()        , 42);
	is_eq( cap()        , 42);
	
	diag("Function pointer inside function")
	is_eq(action(add2), add2(2,3,4));
	is_eq(action(mul2), mul2(2,3,4));

	diag("Function with empty return")
	{
		empty_return();
		pass("ok");
	}
	{
		empty_return_int(2);
		pass("ok");
	}
	{
		int Y = 6;
		is_eq(empty_return_int2(&Y),1);
	}
	{
		double Y = 6;
		is_eq(empty_return_double(Y),1);
	}
	{
		double Y = 6;
		is_eq(empty_return_double2(&Y),1);
	}
	{
		is_eq((empty_return_struct(6)).re,1);
	}
	{
		is_eq((*(empty_return_struct2(6))).re,1);
	}

	diag("typedef function");
	{
		int result = call_a_func(&mul);
		is_eq(result,40);
	}
	{
		is_eq(call_a_func(&mul),40);
	}
	
	diag("function name like in CSTD");
	{
		is_eq(tolower(34,52),86);
		is_eq(toupper(34,52),86);
	}

    done_testing();
}

void my_function()
{
    pass("%s", "Welcome to my function. Feel at home.");
}
