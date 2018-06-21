// Array examples

#include "tests.h"
#include <stdlib.h>

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

double rep_double(double a)
{
	return a;
}

int rep_int(int a)
{
	return a;
}

void zero(int *a, int *b, int *c)
{
	*a = *b = *c = 0;
}

float * next_pointer(float *v)
{
	long l = 1;
	long p = 2;
	(void)(l);
	(void)(p);
	return p - p + v + l;
}

double *dvector(long nl, long nh)
{
	double *v;
	v=(double *)malloc((size_t) ((nh-nl+1+1)*sizeof(double)));
	for (int i=0;i<nh-nl;i++){
		*(v + i) = 42.0;
	}
	return v-nl+1;
}

typedef struct s structs;
void test_pointer_arith_size_t()
{
	size_t size = 1;
	char *left_ptr;
	char arr[3];
	arr[0] = 'a';
	arr[1] = 'b';
	arr[2] = 'c';
	left_ptr = &arr;
	is_eq(*left_ptr , arr[0]);
	left_ptr  = left_ptr + size;
	is_eq(*left_ptr , arr[1]);
	left_ptr += size;
	is_eq(*left_ptr , arr[2]);

	// tests for pointer to struct with size > 1
	structs a[] = {{1, 'a'}, {2, 'b'}, {3, 'c'}};
    is_eq(a[0].i, 1);
    is_eq(a[0].c, 'a');
    is_eq(a[1].i, 2);
    is_eq(a[1].c, 'b');
    is_eq(a[2].i, 3);
    is_eq(a[2].c, 'c');
    structs *ps = &a;
    structs *ps2;
    is_eq(ps->i, 1);
    ps2 = ps + size;
    is_eq(ps2->i, 2);
    ps2 += size;
    is_eq(ps2->i, 3);
    is_eq(ps2-ps, 2);
    ps2 -= size;
    is_eq(ps2->i, 2);
}

void test_pointer_minus_pointer()
{
	char *left_ptr;
	char *right_ptr;
	char arr[30];
	left_ptr  = &arr[0];
	right_ptr = &arr[20];

	is_eq(right_ptr - left_ptr, 20);

	// tests for pointer to struct with size > 1
	structs arr2[30];
	structs *left_ptr2 = &arr2[0];
	structs *right_ptr2 = &arr2[20];

	is_eq(right_ptr2 - left_ptr2, 20);
}

typedef unsigned char pcre_uchar;
typedef unsigned char pcre_uint8;
typedef unsigned short pcre_uint16;
typedef unsigned int pcre_uint32;

#define PT_ANY        0    /* Any property - matches all chars */
#define PT_SC         4    /* Script (e.g. Han) */

#define CHAR_B 'b'
#define STR_A                       "\101"
#define STR_a                       "\141"
#define STR_b                       "\142"
#define STR_c                       "\143"
#define STR_e                       "\145"
#define STR_i                       "\151"
#define STR_m                       "\155"
#define STR_n                       "\156"
#define STR_r                       "\162"
#define STR_y                       "\171"

#define STRING_Any0 STR_A STR_n STR_y "\0"
#define STRING_Arabic0 STR_A STR_r STR_a STR_b STR_i STR_c "\0"
#define STRING_Armenian0 STR_A STR_r STR_m STR_e STR_n STR_i STR_a STR_n "\0"

const char _test_utt_names[] =
  STRING_Any0
  STRING_Arabic0
  STRING_Armenian0;

enum {
  ucp_Arabic,
  ucp_Armenian,
};

typedef struct {
  pcre_uint16 name_offset;
  pcre_uint16 type;
  pcre_uint16 value;
} ucp_type_table;

const ucp_type_table _test_utt[] = {
  {   0, PT_ANY, 0 },
  {   4, PT_SC, ucp_Arabic },
  {  11, PT_SC, ucp_Armenian }
};

int comp (char* name, int i) {
    return strcmp(name, _test_utt_names + _test_utt[i].name_offset);
}

void test_arr_to_pointer() {
    is_true(comp("Any", 0) == 0);
    is_true(comp("Any", 1) != 0);
    is_true(comp("Arabic", 1) == 0);
    is_true(comp("Arabic", 2) != 0);
    is_true(comp("Armenian", 2) == 0);
    is_true(comp("Armenian", 1) != 0);
    pcre_uint32 copynames[1024];
    pcre_uint32 *copynames32 = (pcre_uint32 *)copynames;
    *copynames32 = 42;
    is_eq(copynames[0], 42);
    *copynames = 0;
    is_eq(copynames[0], 0);
    pcre_uint32 c = 0;
    pcre_uchar buffer[8];
    buffer[0] = 7;
    is_eq(c, 0);
    c = *buffer;
    is_eq(c, 7);
}

int main()
{
    plan(162);

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

 	diag("Pointer to Pointer. 1")
 	{
 		double Var = 42;
 		double **PPptr1;
 		double * PPptr2;
 		PPptr2 = &Var;
 		PPptr1 = &PPptr2;
 		is_eq(**PPptr1,Var)
 		Var = 43;
 		is_eq(**PPptr1,Var)
 		(void)(PPptr1);
 		(void)(PPptr2);
 	}
 	diag("Pointer to Pointer. 2")
 	{
 		double Var = 42.0, **PPptr1, * PPptr2;
 		PPptr2 = &Var;
 		PPptr1 = &PPptr2;
 		is_eq(**PPptr1,Var)
 		Var = 43.0;
 		is_eq(**PPptr1,Var)
 		(void)(PPptr1);
 		(void)(PPptr2);
 	}
	diag("Pointer to Pointer. 3");
	{
		int i = 50;
		int ** ptr1;
		int *  ptr2;
		ptr2 = &i;
		ptr1 = &ptr2;
		is_eq(**ptr1, i);
		is_eq(* ptr2, i);
	}
	diag("Pointer to Pointer. 4");
	{
		double arr[5] = {10.,20.,30.,40.,50.};
		double *ptr ;
		ptr = &arr;
		is_eq(*ptr, 10.);
		++ptr;
		is_eq(*ptr, 20.);
	}
	diag("Pointer to Pointer. 5");
	{
		double arr[5] = {10.,20.,30.,40.,50.};
		double *ptr ;
		ptr = &arr;
		is_eq(*ptr, 10.);
		ptr += 1;
		is_eq(*ptr, 20.);
	}
	diag("Pointer to Pointer. 6");
	{
		int arr[5] = {10,20,30,40,50};
		int *ptr ;
		ptr = &arr;
		is_eq(*ptr, 10);
		ptr = 1 + ptr;
		is_eq(*ptr, 20);
	}
	diag("Pointer to Pointer. 7");
	{
		double arr[5] = {10.,20.,30.,40.,50.};
		double *ptr ;
		ptr = &arr;
		is_eq(*ptr, 10.);
		ptr = 1 + ptr;
		is_eq(*ptr, 20.);
	}
	diag("Pointer to Pointer. 8");
	{
		double arr[5] = {10.,20.,30.,40.,50.};
		double *ptr ;
		ptr = &arr;
		is_eq(*ptr, 10.);
		ptr++;
		is_eq(*ptr, 20.);
	}
	diag("Pointer to Pointer. 9");
	{
		double arr[5] = {10.,20.,30.,40.,50.};
		double *ptr ;
		ptr = &arr[2];
		is_eq(*ptr, 30.);
		ptr = ptr -1;
		is_eq(*ptr, 20.);
	}
	diag("Pointer to Pointer. 10");
	{
		double arr[5] = {10.,20.,30.,40.,50.};
		double *ptr ;
		ptr = &arr[2];
		is_eq(*ptr, 30.);
		ptr -= 1;
		is_eq(*ptr, 20.);
	}
	diag("Pointer to Pointer. 11");
	{
		double arr[5] = {10.,20.,30.,40.,50.};
		double *ptr ;
		ptr = &arr[2];
		is_eq(*ptr, 30.);
		ptr--;
		is_eq(*ptr, 20.);
	}
	diag("Pointer to Pointer. 12");
	{
		double arr[5] = {10.,20.,30.,40.,50.};
		double *ptr ;
		int i = 0;
		for (ptr = &arr[0]; i < 5; ptr++){
			is_eq(*ptr,arr[i]);
			i++;
		}
	}
	diag("Operation += 1 for double array");
	{
		float **m;
		m = (float **) malloc(5*sizeof(float*));
		is_not_null(m);
		m[0] = (float *) malloc(10*sizeof(float));
		m[1] = (float *) malloc(10*sizeof(float));
		m[0] += 1;
		(void)(m);
		pass("ok");
	}
	diag("*Pointer = 0");
	{
		int a,b,c;
		a = b = c = 10;
		is_eq(a , 10);
		zero(&a,&b,&c);
		is_eq(a , 0);
		is_eq(b , 0);
		is_eq(c , 0);
		pass("ok");
	}
	diag("pointer + long");
	{
		float *v = (float *)malloc(5*sizeof(float));
		*(v+0) = 5;
		*(v+1) = 6;
		is_eq(*(next_pointer(v)),6);
	}
	diag("create array");
	{
		double * arr = dvector(1,12);
		is_not_null(arr);
		is_eq(arr[1],42.0);
		is_eq(arr[9],42.0);
		(void)(arr);
	}

	diag("Increment inside array 1");
	{
		float f[4] = {1.2,2.3,3.4,4.5};
		int iter = 0;
		is_eq(f[iter++] , 1.2);
		is_eq(f[iter+=1], 3.4);
		is_eq(f[--iter] , 2.3);
	}
	diag("Increment inside array 2");
	{
		struct struct_I_A{
			double * arr;
			int    * pos;
		} ;
		struct struct_I_A siia[2];
		{
			double t_arr [5];
			siia[0].arr = t_arr;
		}
		{
			double t_arr [5];
			siia[1].arr = t_arr;
		}
		{
			int t_pos[1];
			siia[0].pos = t_pos;
		}
		{
			int t_pos[1];
			siia[1].pos = t_pos;
		}
		int t = 0;
		int ii,jj;
		int one = 1;

		siia[0].arr[0] = 45.;
		siia[0].arr[1] = 35.;
		siia[0].arr[2] = 25.;

		siia[0].pos[0] = 0;
		ii = -1;
		jj = -1;
		is_eq(siia[0].arr[(t ++, siia[jj += one].pos[ii += one] += one, siia[jj].pos[ii])] , 35.);

		siia[0].pos[0] = 0;
		ii = -1;
		jj = -1;
		is_eq(siia[0].arr[(t ++, siia[++ jj  ].pos[++ ii ] ++, siia[jj].pos[ii]  )] , 35.);

		siia[0].pos[0] = 2;
		ii = -1;
		jj = -1;
		is_eq(siia[0].arr[(t ++, siia[0].pos[ii += 1] -= 1, siia[0].pos[ii])] , 35.);

		siia[0].pos[0] = 2;
		ii = -1;
		jj = -1;
		is_eq(siia[0].arr[(t ++, siia[0].pos[ii += 1] -- , siia[0].pos[ii] )] , 35.);

		is_eq(t,4);
		(void)(t);
	}
	diag("Increment inside array 3");
	{
		struct struct_I_A3{
			double*arr;
			int    pos;
		} ;
		struct struct_I_A3 siia[2];
		{
			double t_arr [5];
			siia[0].arr = t_arr;
		}
		{
			double t_arr [5];
			siia[1].arr = t_arr;
		}

		siia[0].arr[0] = 45.;
		siia[0].arr[1] = 35.;
		siia[0].arr[2] = 25.;

		siia[0].pos = 0;
		is_eq(siia[0].arr[siia[0].pos += 1] , 35.);

		siia[0].pos = 0;
		is_eq(siia[0].arr[siia[0].pos ++  ] , 45.);

		siia[0].pos = 0;
		is_eq(siia[0].arr[++ siia[0].pos  ] , 35.);

		siia[0].pos = 2;
		is_eq(siia[0].arr[siia[0].pos -= 1] , 35.);

		siia[0].pos = 2;
		is_eq(siia[0].arr[siia[0].pos --  ] , 25.);
	}
	diag("Increment inside array 4");
	{
		struct struct_I_A4{
			double*arr    ;
			int    pos    ;
		} ;
		struct struct_I_A4 siia[2];
		{
			double t_arr [5];
			siia[0].arr = t_arr;
		}
		{
			double t_arr [5];
			siia[1].arr = t_arr;
		}
		int t = 0;

		siia[0].arr[0] = 45.;
		siia[0].arr[1] = 35.;
		siia[0].arr[2] = 25.;

		siia[0].pos = 0;
		is_eq(siia[0].arr[(t ++ , siia[0].pos += 1)] , 35.);

		siia[0].pos = 0;
		is_eq(siia[0].arr[(t ++ ,siia[0].pos ++  )] , 45.);

		siia[0].pos = 2;
		is_eq(siia[0].arr[(t ++ ,siia[0].pos -= 1)] , 35.);

		siia[0].pos = 2;
		is_eq(siia[0].arr[(t ++, siia[0].pos --  )] , 25.);

		is_eq(t,4);
		(void)(t);
	}
	diag("Increment inside array 5");
	{
		struct struct_I_A5{
			double * arr  ;
			int      pos  ;
		} ;
		struct struct_I_A5 siia[2];
		{
			double t_arr [5];
			siia[0].arr = t_arr;
		}
		{
			double t_arr [5];
			siia[1].arr = t_arr;
		}
		int t = 0;

		siia[0].arr[0] = 45.;
		siia[0].arr[1] = 35.;
		siia[0].arr[2] = 25.;

		siia[0].pos = 0;
		is_eq(siia[0].arr[(t ++ , siia[0].pos += 1, siia[0].pos )] , 35.);

		siia[0].pos = 0;
		is_eq(siia[0].arr[(t ++ ,siia[0].pos ++   , siia[0].pos )] , 35.);

		siia[0].pos = 2;
		is_eq(siia[0].arr[(t ++ ,siia[0].pos -=  1, siia[0].pos )] , 35.);

		siia[0].pos = 2;
		is_eq(siia[0].arr[(t ++, siia[0].pos --   , siia[0].pos )] , 35.);

		is_eq(t,4);
		(void)(t);
	}
	diag("Increment inside array 6");
	{
		struct struct_I_A6{
			double * arr ;
			int    * pos ;
		} ;
		struct struct_I_A6 siia[2];
		{
			double t_arr [5];
			siia[0].arr = t_arr;
		}
		{
			double t_arr [5];
			siia[1].arr = t_arr;
		}
		{
			int t_pos[1];
			siia[0].pos = t_pos;
		}
		{
			int t_pos[1];
			siia[1].pos = t_pos;
		}
		int t = 0;

		siia[0].arr[0] = 45.;
		siia[0].arr[1] = 35.;
		siia[0].arr[2] = 25.;

		siia[0].pos[0] = 0;
		is_eq(siia[0].arr[(t ++, siia[0].pos[0] += 1)] , 35.);

		siia[0].pos[0] = 0;
		is_eq(siia[0].arr[(t ++, siia[0].pos[0] ++  )] , 45.);

		siia[0].pos[0] = 2;
		is_eq(siia[0].arr[(t ++, siia[0].pos[0] -= 1)] , 35.);

		siia[0].pos[0] = 2;
		is_eq(siia[0].arr[(t ++, siia[0].pos[0] --  )] , 25.);

		is_eq(t,4);
		(void)(t);
	}

	test_pointer_arith_size_t();
	test_pointer_minus_pointer();

	diag("negative array index");
    {
        pcre_uchar arr[] = "abcdef";
        pcre_uchar *a = &arr[2];
        is_eq(*a, 'c');
        is_eq(a[-1], 'b');
        is_eq(a[-2+1], 'b');
        is_eq(*(a-1), CHAR_B);
    }

    diag("array to pointer");
    test_arr_to_pointer();

    done_testing();
}
