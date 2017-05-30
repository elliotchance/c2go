#include <stdio.h>
#include "tests.h"

int main()
{
    plan(5);
	
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

	done_testing();
}
