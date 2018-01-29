// Tests for unions.

#include <stdio.h>
#include "tests.h"

//////////////////////////////
//
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

//////////////////////////////
//
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

//////////////////////////////
//
struct SHA3 {
  union {
    double iY;
    double dY;
  } uY;
  float ffY;
};

union unknown {
  double i2;
  double d2;
};
struct SHA32 {
  union unknown u2;
  float ff2;
};


void union_in_struct()
{
	diag("Union inside struct");
	struct SHA3 sha;
	sha.ffY  = 12.444;
	sha.uY.iY = 4;
	is_eq(sha.uY.iY, 4);
	is_eq(sha.uY.dY, 4);
	is_eq(sha.ffY , 12.444);

	struct SHA32 sha2;
	sha2.ff2  = 12.444;
	sha2.u2.i2 = 4;
	is_eq(sha2.u2.i2, 4);
	is_eq(sha2.u2.d2, 4);
	is_eq(sha2.ff2 , 12.444);
}

//////////////////////////////
//
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

//////////////////////////////
//

typedef int ii;
typedef struct SHA SHA;
struct SHA{
  union {
    ii            s[25];
    unsigned char x[100];
  } u;
  unsigned uuu;
};

void union_arr_in_str()
{
	SHA sha;
	sha.uuu = 15;
	is_eq(sha.uuu,15);
	for (int i = 0 ; i< 25;i++)
		sha.u.s[0] = 0;
	is_eq(sha.u.s[0],0);
	is_true(sha.u.x[0] == 0);
	for (int i=0;i<6;i++){
		sha.u.s[i] = (ii)(4);
		sha.u.s[i] = (ii)(42) + sha.u.s[i];
	}
	is_eq(sha.u.s[5],46);
	is_true(sha.u.x[0] != 0);
}

//////////////////////////////
//
int main()
{
    plan(30);

	union_simple     ();
	union_array      ();
	union_in_struct  ();
	union_typedef    ();
	union_arr_in_str ();

    done_testing();
}
