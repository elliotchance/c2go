#include <stdio.h>
#include "tests.h"

int main()
{
    plan(15);
	
	int i = 0;

	// There will be 4 checks in the first loop.
	do {
		pass("%s", "first do statement");
		i = i + 1;
	} while( i < 4 );

	// Only one check in the second loop.
	i = 0;
	do {
		i++;
		if(i < 3) continue;
		pass("%s", "second do statement");
	} while(i < 3);

	diag("check while");
	i = 1000;
	do {
		i--;
		if (i < 10) { break; }
	} while ((i /= 10) > 0);
	is_eq( i , 8 );

	diag("do without CompoundStmt");
	int s = 1;
	do s++; while(s < 10);
	is_eq(s , 10);

	diag("do while non-boolean expression")
	{
        int n = 3;
        do {
            pass("%d", n);
        } while(--n);
        is_eq(n, 0)
	}

	{
        int n = -3;
        do {
            pass("%d", n);
        } while(++n);
        is_eq(n, 0)
	}

	done_testing();
}
