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

typedef unsigned char pcre_uchar;
typedef unsigned char pcre_uint8;
typedef struct pcre_study_data {
    pcre_uint8 start_bits[32];
} pcre_study_data;

void caststr() {
    pcre_uchar str[] = "abcd";
    is_streq((char *) str, "abcd");
}

static const pcre_uchar TEST[] =  {
  'x', (pcre_uchar) CHAR_NBSP, '\n', '\0' };

#define CHAR_E ((unsigned char) 'e')
static const pcre_uchar TEST2[] =  {
  'x', (pcre_uchar) CHAR_NBSP, '\n',
  (pcre_uchar) CHAR_E, '\0' };


void test_static_array()
{
    is_eq('x', TEST[0]);
    is_eq((pcre_uchar) '\xa0', TEST[1]);
    is_eq('\n', TEST[2]);
    is_eq('x', TEST2[0]);
    is_eq('e', TEST2[3]); // can distinguish character at same column in different lines
}

void castbitwise() {
    pcre_uint32 x = 0xff;
    x &= ~0x3c;
    is_eq(x, 0xc3);
}

void cast_pointer_diff(pcre_uchar *str, int *x) {
    pcre_uchar *p = str;
    pcre_uchar ab = '\0';
    *x = (int)(p - str) - ab;
}

typedef unsigned char x_uint;
typedef unsigned char y_uint;
typedef struct {
    unsigned char a;
    unsigned char b;
} z_struct;

void test_voidcast()
{
    x_uint x = 42;
    void * y = &x;
    y_uint *z = (y_uint*) y;
    is_eq(42, *z);
    x_uint arr1[] = { 1, 2, 3, 4 };
    y = arr1;
    z_struct *arr2 = (z_struct*) y;
    is_eq(1, arr2[0].a);
    is_eq(2, arr2[0].b);
    is_eq(3, arr2[1].a);
    is_eq(4, arr2[1].b);

    pcre_uchar **stringlist;
    y = &x;
    void * py = &y;
    stringlist = (pcre_uchar **)(py);
    is_eq(**((x_uint **)py), 42);
    is_eq(**stringlist, 42);

    *((const char **)py) = NULL;
    is_eq(x, 42);
    is_true(y == NULL);
    y = &x;
    *((const pcre_uint8 **)py) = NULL;
    is_eq(x, 42);
    is_true(y == NULL);
    y = &x;
    *((const pcre_uchar **)py) = NULL;
    is_eq(x, 42);
    is_true(y == NULL);
}

int main()
{
    plan(52);

    START_TEST(cast);
    START_TEST(castbool);
    START_TEST(vertex);
    START_TEST(strCh);
    START_TEST(voidcast);

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

    diag("Typedef slice convertion")
    caststr();

    diag("Compare with static array")
    test_static_array();

    diag("Cast with compound assign operator")
    castbitwise();

    diag("Cast pointer diff");
    {
        pcre_uchar s[] = "abcd";
        int b = 42;
        cast_pointer_diff(&s[0], &b);
        is_eq(b, 0);
    }
	diag("Cast array to slice");
    {
        pcre_study_data sdata;
        sdata.start_bits[1] = 42;
        const pcre_study_data *study = &sdata;
        const pcre_uint8 *p = 0;
        p = study->start_bits;
        is_eq(p[1], 42);
    }

    done_testing();
}
