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
	plan(30);

	// Amount test cases
	int amount_cases = 10;
	// Test cases
	testCase tcs[10];
	// case 0
	tcs[0].argc = 0;
	tcs[0].argv = NULL;
	tcs[0].aflag = 0;
	tcs[0].bflag = 0;
	tcs[0].cvalue = NULL;

	// case 1
	tcs[1].argc = 3;
	{
		char *v0 = "programName";
		char *v1 = "-a";
		char *v2 = "-b";
		char *c[3];
		c[0] = v0;
		c[1] = v1;
		c[2] = v2;
		tcs[1].argv = (char**)c;
	}
	tcs[1].aflag = 1;
	tcs[1].bflag = 1;
	tcs[1].cvalue = NULL;

	// case 2
	tcs[2].argc = 2;
	{
		char *v0 = "programName";
		char *v1 = "-ab";
		char *c[2];
		c[0] = v0;
		c[1] = v1;
		tcs[2].argv = (char**)c;
	}
	tcs[2].aflag = 1;
	tcs[2].bflag = 1;
	tcs[2].cvalue = NULL;

	// case 3
	tcs[3].argc = 3;
	{
		char *v0 = "programName";
		char *v1 = "-c";
		char *v2 = "foo";
		char *c[3];
		c[0] = v0;
		c[1] = v1;
		c[2] = v2;
		tcs[3].argv = (char**)c;
	}
	tcs[3].aflag = 0;
	tcs[3].bflag = 0;
	tcs[3].cvalue = "foo";

	// case 4
	tcs[4].argc = 2;
	{
		char *v0 = "programName";
		char *v1 = "-cfoo";
		char *c[2];
		c[0] = v0;
		c[1] = v1;
		tcs[4].argv = (char**)c;
	}
	tcs[4].aflag = 0;
	tcs[4].bflag = 0;
	tcs[4].cvalue = "foo";

	// case 5
	tcs[5].argc = 2;
	{
		char *v0 = "programName";
		char *v1 = "arg1";
		char *c[2];
		c[0] = v0;
		c[1] = v1;
		tcs[5].argv = (char**)c;
	}
	tcs[5].aflag = 0;
	tcs[5].bflag = 0;
	tcs[5].cvalue = NULL;

	// case 6
	tcs[6].argc = 3;
	{
		char *v0 = "programName";
		char *v1 = "-a";
		char *v2 = "arg1";
		char *c[3];
		c[0] = v0;
		c[1] = v1;
		c[2] = v2;
		tcs[6].argv = (char**)c;
	}
	tcs[6].aflag = 1;
	tcs[6].bflag = 0;
	tcs[6].cvalue = NULL;

	// case 7
	tcs[7].argc = 4;
	{
		char *v0 = "programName";
		char *v1 = "-c";
		char *v2 = "foo";
		char *v3 = "arg1";
		char *c[4];
		c[0] = v0;
		c[1] = v1;
		c[2] = v2;
		c[3] = v3;
		tcs[7].argv = (char**)c;
	}
	tcs[7].aflag = 0;
	tcs[7].bflag = 0;
	tcs[7].cvalue = "foo";

	// case 8
	tcs[8].argc = 4;
	{
		char *v0 = "programName";
		char *v1 = "-a";
		char *v2 = "--";
		char *v3 = "-b";
		char *c[4];
		c[0] = v0;
		c[1] = v1;
		c[2] = v2;
		c[3] = v3;
		tcs[8].argv = (char**)c;
	}
	tcs[8].aflag = 1;
	tcs[8].bflag = 0;
	tcs[8].cvalue = NULL;

	// case 9
	tcs[9].argc = 3;
	{
		char *v0 = "programName";
		char *v1 = "-a";
		char *v2 = "-";
		char *c[3];
		c[0] = v0;
		c[1] = v1;
		c[2] = v2;
		tcs[9].argv = (char**)c;
	}
	tcs[9].aflag = 1;
	tcs[9].bflag = 0;
	tcs[9].cvalue = NULL;

	int i;
	for ( i = 0; i < amount_cases ; i++)
	{
		diag("Test case");
		int aflag, bflag;
		aflag = 0;
		bflag = 0;
		char *cvalue = NULL;
		int c = 0;

		opterr = 1;
		optind = 1;
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
		if (tcs[i].cvalue == NULL && cvalue == NULL) 
		{
			pass("both is nil");
		}
		else
		{
			if (tcs[i].cvalue != NULL && cvalue == NULL) {
				fail("fail cvalue is nil");
			} else if (tcs[i].cvalue == NULL && cvalue != NULL) {
				fail("fail cvalue is not nil");
			} else {
				is_streq( tcs[i].cvalue, cvalue );
			}
		}
	}

    done_testing();
}
