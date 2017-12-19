// Array examples

#include "tests.h"

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_intarr()
{
    int a[3];
    a[0] = 5;
    a[1] = 9;
    a[2] = -13;

    is_eq(a[0], 5);
    is_eq(a[1], 9);
    is_eq(a[2], -13);
}

void test_doublearr()
{
    double a[2];
    a[0] = 1.2;
    a[1] = 7; // different type

    is_eq(a[0], 1.2);
    is_eq(a[1], 7.0);
}

void test_intarr_init()
{
    int a[] = {10, 20, 30};
    is_eq(a[0], 10);
    is_eq(a[1], 20);
    is_eq(a[2], 30);
}

void test_floatarr_init()
{
    float a[] = {2.2, 3.3, 4.4};
    is_eq(a[0], 2.2);
    is_eq(a[1], 3.3);
    is_eq(a[2], 4.4);
}

void test_chararr_init()
{
    char a[] = {97, 98, 99};
    is_eq(a[0], 'a');
    is_eq(a[1], 'b');
    is_eq(a[2], 'c');
}

void test_chararr_init2()
{
    char a[] = {'a', 'b', 'c'};
    is_eq(a[0], 'a');
    is_eq(a[1], 'b');
    is_eq(a[2], 'c');
}

void test_exprarr()
{
    int a[] = {2 ^ 1, 3 & 1, 4 | 1, (5 + 1)/2};
    is_eq(a[0], 3);
    is_eq(a[1], 1);
    is_eq(a[2], 5);
    is_eq(a[3], 3);
}

struct s {
    int i;
    char c;
};

void test_structarr()
{
    struct s a[] = {{1, 'a'}, {2, 'b'}};
    is_eq(a[0].i, 1);
    is_eq(a[0].c, 'a');
    is_eq(a[1].i, 2);
    is_eq(a[1].c, 'b');

    struct s b[] = {(struct s){1, 'a'}, (struct s){2, 'b'}};
    is_eq(b[0].i, 1);
    is_eq(b[0].c, 'a');
    is_eq(b[1].i, 2);
    is_eq(b[1].c, 'b');
}

long dummy(char foo[42])
{
    return sizeof(foo);
}

void test_argarr()
{
    char abc[1];
    is_eq(8, dummy(abc));
}

void test_multidim()
{
    int a[2][3] = {{5,6,7},{50,60,70}};
    is_eq(a[1][2], 70);

    // omit array length
    int b[][3][2] = {{{1,2},{3,4},{5,6}},
                     {{6,5},{4,3},{2,1}}};
    is_eq(b[1][1][0], 4);
    // 2 * 3 * 2 * sizeof(int32)
    is_eq(sizeof(b), 48);

    struct s c[2][3] = {{{1,'a'},{2,'b'},{3,'c'}}, {{4,'d'},{5,'e'},{6,'f'}}};
    is_eq(c[1][1].i, 5);
    is_eq(c[1][1].c, 'e');
    c[1][1] = c[0][0];
    is_eq(c[1][1].i, 1);
    is_eq(c[1][1].c, 'a');
}

void test_ptrarr()
{
    int b = 22;

    int *d[3];
    d[1] = &b;
    is_eq(*(d[1]), 22);

    int **e[4];
    e[0] = d;
    is_eq(*(e[0][1]), 22);
}

void test_stringarr_init()
{
    char *a[] = {"a", "bc", "def"};
    is_streq(a[0], "a");
    is_streq(a[1], "bc");
    is_streq(a[2], "def");
}

void test_partialarr_init()
{
    // Last 2 values are filled with zeros
    double a[4] = {1.1, 2.2};
    is_eq(a[2], 0.0);
    is_eq(a[3], 0.0);

    struct s b[3] = {{97, 'a'}};
    is_eq(b[0].i, 97);
    is_eq(b[2].i, 0);
    is_eq(b[2].c, 0);
}

extern int arrayEx[];
int arrayEx[4] = { 1, 2, 3, 4 };

int ff(){ return 3;}

int main()
{
    plan(66);

    START_TEST(intarr);
    START_TEST(doublearr);
    START_TEST(intarr_init);
    START_TEST(floatarr_init);
    START_TEST(chararr_init);
    START_TEST(chararr_init2);
    START_TEST(exprarr);
    START_TEST(structarr);
    START_TEST(argarr);
    START_TEST(multidim);
    START_TEST(ptrarr);
    START_TEST(stringarr_init);
    START_TEST(partialarr_init);

	is_eq(arrayEx[1],2.0);

	diag("Array arithmetic")
    float a[5];
    a[0] = 42.;
       is_eq(a[0],42.);
    a[0+1] = 42.;
       is_eq(a[1],42);
    a[2]   = 42.;
       is_eq(a[2],42);
       
    diag("Pointer arithmetic. Part 1");
    float *b;
    b = (float *)calloc(5,sizeof(float));
    
    *b   = 42.;
    is_eq(*(b+0),42.);
    
    *(b+1) = 42.;
    is_eq(*(b+1),42.);
    *(2+b) = 42.;
    is_eq(*(b+2),42.);

    *(b+ff()) = 45.;
    is_eq(*(b + 3), 45.);
    *(ff()+b+1) = 46.;
    is_eq(*(b + 4), 46.);

	*(b+ (0 ? 1 : 2)) = -1.;
	is_eq(*(b+2),-1);

	*(b + 0) = 1 ;
	*(b + (int)(*(b + 0)) - 1) = 35;
	is_eq(*(b+0),35);

	*(b + (int)((float)(2))) = -45;
	is_eq(*(b+2),-45);

	*(b + 1 + 3 + 1 - 5*1 + ff() - 3) = -4.0;
	is_eq(*(b+0), -4.0);
	is_eq(*b    , -4.0);

	is_eq((*(b + 1 + 3 + 1 - 5*1 + ff() - 3 + 1) = -48.0,*(b+1)), -48.0);
	{int rrr;(void)(rrr);}
	
	diag("Pointer arithmetic. Part 2")
	{
		float *arr; 
		arr = (float*)calloc(1+1,sizeof(float)); 
		is_true(arr != NULL);
		(void)(arr);
	}
	{
		float *arr;
		arr = (float *) calloc(1+ff(),sizeof(float));
		is_true(arr != NULL);
		(void)(arr);
	}
	{
		float *arr;
		arr = (float *) calloc(ff()+ff(),sizeof(float));
		is_true(arr != NULL);
		(void)(arr);
	}
	{
		float *arr;
		arr = (float *) calloc(ff()+1+0+0+1*0,sizeof(float));
		is_true(arr != NULL);
		(void)(arr);
	}

    done_testing();
}
