#include <stdio.h>
#include "tests.h"

int d(int v){
	return 2*v;
}

int main()
{
    plan(7);

    int x = 1;

    // Without else
    if (x == 1)
        pass("%s", "x is equal to one");
	
    if (x == 1){
        pass("%s", "x is equal to one");
	}

    // With else
    if (x != 1){
        fail("%s", "x is not equal to one");
	} else {
        pass("%s", "x is equal to one");
	}

	if ( NULL) {
		fail("%s", "NULL is zero");
	} else {
		pass("%s", "NULL is not zero");
	}

	if ( ! NULL) {
		pass("%s", "Invert : ! NULL is zero");
	} else {
		fail("%s", "Invert : ! NULL is not zero");
	}

	diag("Operation inside function")
	int ii = 5;
	if ((ii = d(ii)) != (-1)){
		is_eq(ii,10)
	}
	
	diag("if - else");
	{
		int a = 10;
		int b = 5 ;
		int c = 0 ;
		if ( a < b )
		{
		} else if (b > c)
		{
			pass("ok");
		}
	}

    done_testing();
}
