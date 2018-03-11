#include <stdio.h>
#include "tests.h"

typedef float **triangle;
#define deadtri(tria)  ((tria)[1] == (triangle) NULL)

int main()
{
    plan(8);

    int value = 1;

    while (value <= 3)
    {
        pass("value is %d", value);
        value++;
    }

    // continue
    value = 0;
    while (value < 3)
    {
        value++;
        if (value < 3)
            continue;
        pass("%d", value);
    }

	diag("while without body")
	while(0);
	pass("%s","while without body");

	value = 1;
	while((value--,value));
	is_eq(value , 0);

	diag("while with star");
	{
	int * ok;
	int value2;
	ok = & value2;
	*ok = 1;
	int iterator = 0;
	do{
		if (iterator == 1){
			*ok = 0;
		}
		iterator ++;
		if (iterator >10){
			break;
		}
	}while(*ok);
	is_eq(*ok, 0);
	}

	diag("while with --");
	{
		int T = 2;
		int counter = 0;
		while(T--){
			if (counter > 50){
				break;
			}
		};
		is_eq(T,-1);
	}

    done_testing();
}
