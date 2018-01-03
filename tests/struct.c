// Tests for structures.

#include <stdio.h>
#include "tests.h"

struct programming
{
    float constant;
    char *pointer;
};

void pass_by_ref(struct programming *addr)
{
    char *s = "Show string member.";
    float v = 1.23+4.56;

    addr->constant += 4.56;
    addr->pointer = s;

    is_eq(addr->constant, v);
    is_streq(addr->pointer, "Show string member.");
}

void pass_by_val(struct programming value)
{
    value.constant++;

    is_eq(value.constant, 2.23);
    is_streq(value.pointer, "Programming in Software Development.");
}

typedef struct mainStruct{
    double constant;
} secondStruct;

typedef struct {
    double t;
} ts_c;

typedef struct ff {
    int v1,v2;
} tt1, tt2;

struct outer {
    int i;
    struct z {
        int j;
    } inner;
};

struct xx {
    int i;
    struct yy {
        int j;
        struct zz {
            int k;
        } deep;
    } inner;
};

int summator(int i, float f){
	return i+(int)(f);
}

int main()
{
    plan(46);

    struct programming variable;
    char *s = "Programming in Software Development.";

    variable.constant = 1.23;
    variable.pointer = s;

    is_eq(variable.constant, 1.23);
    is_streq(variable.pointer, "Programming in Software Development.");

    pass_by_val(variable);
    pass_by_ref(&variable);

    struct mainStruct s1;
    s1.constant = 42.;
    is_eq(s1.constant, 42.);

    secondStruct s2;
    s2.constant = 42.;
    is_eq(s2.constant, 42.);

    ts_c c;
    c.t = 42.;
    is_eq(c.t , 42.);

    tt1 t1;
    t1.v1 = 42.;
    is_eq(t1.v1,42.)

    tt2 t2;
    t2.v1 = 42.;
    is_eq(t2.v1,42.)

    struct ff tf2;
    tf2.v2 = t1.v1;
    is_eq(tf2.v2,t1.v1);

    struct outer o;
    o.i = 12;
    o.inner.j = 34;
    is_eq(o.i + o.inner.j, 46);

    struct xx x;
    x.i = 12;
    x.inner.j = 34;
    x.inner.deep.k = 56;
    is_eq(x.i + x.inner.j + x.inner.deep.k, 102);

	struct u{
		int y;
	};
	struct u yy;
	yy.y = 42;
	is_eq(yy.y,42);
	
	diag("Typedef struct with same name")
	{
		typedef struct Uq Uq;
		struct Uq{
			int uq;
		};
		Uq uu;
		uu.uq = 42;
		is_eq(uu.uq,42);
	}

	diag("Initialization of struct")
	struct Point {
		int x;
		int y;
	};
	struct Point p = { .y = 2, .x = 3 };
	is_eq(p.x, 3);
	is_eq(p.y, 2);

	diag("ImplicitValueInitExpr")
	{
		typedef struct {
		    int x2;
		    int y2;
		} coord2;

		typedef struct {
		    coord2 position2;
		    int possibleSteps2;
		} extCoord2;

		extCoord2 followingSteps[2] =
	    {
	        {.possibleSteps2 = 1}, {.possibleSteps2 = 1},
	    };
		is_eq(followingSteps[0].possibleSteps2, 1);
	}
	{
		struct coord{
		    int x;
		    int y;
		};

		struct extCoord{
		    struct coord position;
		    int possibleSteps;
		};

		struct extCoord followingSteps[2] =
	    {
	        {.possibleSteps = 1}, {.possibleSteps = 1},
	    };
		is_eq(followingSteps[0].possibleSteps, 1);
	}

	diag("Double typedef type")
	{
		typedef int  int2;
		typedef int2 int3;
		typedef int3 int4;

		is_eq((int)((int4)((int3)((int2)(42)))),42);
	}
	{
		typedef size_t size2;
		is_eq(((size2)((size_t)(56))),56.0)
	}
	{
		is_eq((size_t)(43),43);
	}

	diag("Function pointer inside struct")
	{
		struct F1{
			  int x;
			  int (*f)(int, float);
		};
		struct F1 f1;
		f1.x = 42;
		f1.f = summator;
		is_eq(f1.x,42);
		is_eq(f1.f(3,5),8);
	}
	{
		typedef struct {
			  int x;
			  int (*f)(int, float);
		} F2;
		F2 f2;
		f2.x = 42;
		f2.f = summator;
		is_eq(f2.x,42);
		is_eq(f2.f(3,5),8);
	}

	diag("typedef function")
	{
		typedef int ALIAS (int, float);
		ALIAS * f = summator;
		is_eq(f(3,5),8);
	}
	{
		typedef int ALIAS2 (int, float);
		ALIAS2 * f;
		f = summator;
		is_eq(f(3,5),8);
	}

	diag("typedef struct C C inside function")
	{
		typedef struct CCC CCC;
		struct CCC {
			float ff;
		};
		CCC c;
		c.ff = 3.14;
		is_eq(c.ff,3.14);
	}
	typedef struct CP CP;
	struct CP {
		float ff;
	};
	CP cp;
	cp.ff = 3.14;
	is_eq(cp.ff,3.14);

	diag("struct name from Go keyword")
	{ struct chan        {int i;}; struct chan        UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct defer       {int i;}; struct defer       UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct fallthrough {int i;}; struct fallthrough UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct func        {int i;}; struct func        UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct go          {int i;}; struct go          UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct import      {int i;}; struct import      UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct interface   {int i;}; struct interface   UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct map         {int i;}; struct map         UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct package     {int i;}; struct package     UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct range       {int i;}; struct range       UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct select      {int i;}; struct select      UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct type        {int i;}; struct type        UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct var         {int i;}; struct var         UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct _           {int i;}; struct _           UU;  UU.i = 5; is_eq(UU.i,5);}
	{ struct init        {int i;}; struct init        UU;  UU.i = 5; is_eq(UU.i,5);}

	// uncomment after success implementation of struct scope
	// https://github.com/elliotchance/c2go/issues/368
/*
	diag("Typedef struct name from Go keyword")
	{ typedef struct {int i;} chan        ;	chan        UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} defer       ;	defer       UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} fallthrough ;	fallthrough UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} func        ;	func        UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} go          ;	go          UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} import      ;	import      UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} interface   ;	interface   UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} map         ;	map         UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} package     ;	package     UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} range       ;	range       UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} select      ;	select      UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} type        ;	type        UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} var         ;	var         UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} _           ;	_           UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} init        ;	init        UU; UU.i = 5; is_eq(UU.i,5);}
*/

    done_testing();
}
