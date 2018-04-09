#include <stdio.h>
#include "tests.h"

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_cast()
{
    int c[] = {(int)'a', (int)'b'};
    is_eq(c[0], 97);

    double d = (double) 1;
    is_eq(d, 1.0);
}

void test_castbool()
{
    char i1 = (1 == 1);
    short i2 = (1 == 1);
    int i3 = (1 == 1);
    long i4 = (1 == 0);
    long long i5 = (1 == 0);

    is_eq((i1==1) && (i2==1) && (i3==1) && (i4==0) && (i5==0), 1);
}

void char_overflow()
{
	{
	char c;	c = -1;
	unsigned char u = c;
	is_eq(u, 256-1);
	}
	{
	char c = -1;
	unsigned char u = c;
	is_eq(u, 256-1);
	}
	{
	char c = (-1);
	unsigned char u = c;
	is_eq(u, 256-1);
	}
	{
	char c = (((-1)));
	unsigned char u = c;
	is_eq(u, 256-1);
	}
}

typedef double * vertex;
void test_vertex()
{
	diag("vertex");

	double a[1];
	a[0] = 42;
	double b[1];
	b[0] = 45;

	double dxoa;
	vertex triorg  = (vertex)(a);
	vertex triapex = (vertex)(b);
	dxoa = triorg[0] - triapex[0];

	is_eq(dxoa, -3);
}

static int strlenChar(const char *z){
  int n = 0;
  while( *z ){
    if( (0xc0&*(z++))!=0x80 ) n++;
  }
  return n;
}

void test_strCh()
{
	char * z = "Hello, c2go\0";
	is_eq(strlenChar(z),11);
}

typedef unsigned int pcre_uint32;
#define CHAR_NBSP                   ((unsigned char)'\xa0')

void test_preprocessor()
{
    int tmp = 160;
    pcre_uint32 chr = tmp;

    is_eq(chr, CHAR_NBSP);
}

int main()
{
    plan(30);

    START_TEST(cast);
    START_TEST(castbool);
    START_TEST(vertex);
    START_TEST(strCh);

	{
	typedef unsigned int u32;
	u32 x = 42;
	is_eq(x , 42);
    u32 a[10];
    a[0] = x;
	is_eq(a[0],42);
	}

	{
	typedef double u32d;
	u32d x = 42.0;
	is_eq(x , 42.0);
    u32d a[10];
    a[0] = x;
	is_eq(a[0],42.0);
	}

	{
	typedef int integer;
	typedef int INTEGER;
    integer i = 123;
    INTEGER j = 567;
    j = i;
    i = j;
	is_eq(i , 123);
	is_eq(j , 123);
	}

	double *d = (double *) 0;
	is_true(d == NULL);
	int    *i = (int    *) 0;
	is_true(i == NULL);
	float  *f = (float  *) 0;
	is_true(f == NULL);
	char   *c = (char   *) 0;
	is_true(c == NULL);

	double *d2 = 0;
	is_true(d2 == NULL);
	int    *i2 = 0;
	is_true(i2 == NULL);
	float  *f2 = 0;
	is_true(f2 == NULL);
	char   *c2 = 0;
	is_true(c2 == NULL);

	diag("Calloc with type")
	{
		double *ddd = (double *)calloc(2,sizeof(double));
		is_not_null(ddd);
		(void)(ddd);
	}
	{
		double *ddd;
		ddd = (double *)calloc(2,sizeof(double));
		is_not_null(ddd);
		(void)(ddd);
	}
	
	diag("Type convertion from void* to ...")
	{
		void * ptr2;
		int tInt = 55;
		ptr2 = &tInt;
		is_eq(*(int*)ptr2, 55);
		double tDouble = -13;
		ptr2 = &tDouble;
		is_eq(*(double*)ptr2,-13);
		float tFloat = 67;
		is_eq(*(float *)(&tFloat),67);
	}
	diag("Type convertion from void* to ... in initialization")
	{
		long tLong = 556;
		void * ptr3 = &tLong;
		is_eq(*(long *) ptr3, 556);
	}

	char_overflow();

    diag("Compare preprocessor with type")
    test_preprocessor();

    done_testing();
}
