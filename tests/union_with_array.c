// Tests for unions.

#include <stdio.h>
#include "tests.h"


typedef int ii;
typedef struct SHA SHA;
struct SHA{
  union {
    ii            s[25];
    unsigned char x[100];
  } u;
  unsigned uuu;
};

void union_with_array()
{
	SHA sha;
	sha.uuu = 15;
	is_eq(sha.uuu,15);
	/* for (int i = 0 ; i< 25;i++) */
		sha.u.s[0] = 0;
	is_eq(sha.u.s[0],0);
	is_true(sha.u.x[0] == 0);
	/* for (int i=0;i<6;i++){ */
	/* 	sha.u.s[i] = (ii)(4); */
	/* 	sha.u.s[i] = (ii)(42) + sha.u.s[i]; */
	/* } */
	/* is_eq(sha.u.s[5],46); */
	/* is_true(sha.u.x[0] != 0); */
}

int main()
{
    plan(3);
	
	union_with_array();
    
	done_testing();
}
