// Test created in according to
// https://www.gnu.org/software/libc/manual/html_node/Example-of-Getopt.html

#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include "tests.h"

typedef struct testCase testCase;
struct testCase
{
	// input data
	int    argc;
	char **argv;
	// output data
	int    aflag;
	int    bflag;
	char *cvalue;
};

int main()
{
    plan(4);

	// Amount test cases
	int amount_cases = 10;
	// Test cases
	testCase tcs[1];
	// case 0
	tcs[0].argc = 0;
	tcs[0].argv = NULL;
	tcs[0].aflag = 0;
	tcs[0].bflag = 0;
	tcs[0].cvalue = NULL;
	// case 1
	tcs[1].argc = 2;
	{
	char out[2][10] = {"-a" , "-b"};
	tcs[1].argv = &out;
	}
	tcs[1].aflag = 1;
	tcs[1].bflag = 1;
	tcs[1].cvalue = NULL;
	// case 2
	tcs[2].argc = 1;
	{
	char out[1][10] = {"-ab"};
	tcs[2].argv = &out;
	}
	tcs[2].aflag = 1;
	tcs[2].bflag = 1;
	tcs[2].cvalue = NULL;
	// case 3
	tcs[3].argc = 2;
	{
	char out[2][10] = {"-c","foo"};
	tcs[3].argv = &out;
	}
	tcs[3].aflag = 0;
	tcs[3].bflag = 0;
	tcs[3].cvalue = "foo";
	// case 4
	tcs[4].argc = 1;
	{
	char out[1][10] = {"-cfoo"};
	tcs[4].argv = &out;
	}
	tcs[4].aflag = 0;
	tcs[4].bflag = 0;
	tcs[4].cvalue = "foo";
	// case 5
	tcs[5].argc = 1;
	{
	char out[1][10] = {"arg1"};
	tcs[5].argv = &out;
	}
	tcs[5].aflag = 0;
	tcs[5].bflag = 0;
	tcs[5].cvalue = NULL;
	// case 6
	tcs[6].argc = 2;
	{
	char out[2][10] = {"-a","arg1"};
	tcs[6].argv = &out;
	}
	tcs[6].aflag = 1;
	tcs[6].bflag = 0;
	tcs[6].cvalue = NULL;
	// case 7
	tcs[7].argc = 3;
	{
	char out[3][10] = {"-c","foo","arg1"};
	tcs[7].argv = &out;
	}
	tcs[7].aflag = 0;
	tcs[7].bflag = 0;
	tcs[7].cvalue = "foo";
	// case 8
	tcs[8].argc = 3;
	{
	char out[3][10] = {"-a","--","-b"};
	tcs[8].argv = &out;
	}
	tcs[8].aflag = 1;
	tcs[8].bflag = 0;
	tcs[8].cvalue = NULL;
	// case 9
	tcs[9].argc = 2;
	{
	char out[2][10] = {"-a","-"};
	tcs[9].argv = &out;
	}
	tcs[9].aflag = 1;
	tcs[9].bflag = 0;
	tcs[9].cvalue = NULL;

	int i;
	for ( i = 0; i < amount_cases ; i++)
	{
		diag("Test case");
		int aflag, bflag;
		char *cvalue;
		int index;
		int c;
		opterr = 0;
		while ((c = getopt (tcs[i].argc, tcs[i].argv, "abc:")) != -1)
			switch (c)
			{
			case 'a': aflag = 1; break;
			case 'b': bflag = 1; break;
			case 'c': cvalue = optarg; break;
			}
		// compare results
		is_eq(tcs[i].aflag , aflag)
		is_eq(tcs[i].bflag , bflag)
		is_streq( tcs[i].cvalue, cvalue );
	}

    done_testing();
}
