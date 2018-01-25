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
	int amount_cases = 1;
	// Test cases
	testCase tcs[1];
	// case 0
	tcs[0].argc = 0;
	tcs[0].argv = NULL;
	tcs[0].aflag = 0;
	tcs[0].bflag = 0;
	tcs[0].cvalue = NULL;

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
