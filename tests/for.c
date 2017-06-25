#include <stdio.h>
#include "tests.h"

int main()
{
    plan(23);

    int i = 0;

    diag("Missing init");
    for (; i < 3; i++)
        pass("%d", i);

    diag("CompountStmt");
    for (i = 0; i < 3; i++)
    {
        pass("%d", i);
    }

    diag("Not CompountStmt");
    for (i = 0; i < 3; i++)
        pass("%d", i);

    diag("Infinite loop");
    int j = 0;
    for (;;)
    {
        pass("%d", j);
        j++;
        if (j > 3)
            break;
    }

    diag("continue");
    i = 0;
    j = 0;
    for (;;)
    {
        pass("%d %d", i, j);
        i++;
        if (i < 3)
            continue;
        j++;
        if (j > 3)
            break;
    }

	diag("big ininitialization");
	for ( i = 0, j = 0 ; i < 2 ; i++){
		pass("%d %d", i, j);
	}
	
	diag("big ininitialization");
	for ( i = 0, j = 0 ; i < 2 ; ){
		pass("%d %d", i, j);
		i++;
		j++;
	}
/*

	diag("big ininitialization");
	int k;
	for ( i = 0, j = 0 , k = 0; i < 2 ; ){
		pass("%d %d %d", i, j, k);
		i++;
		j+=2;
		k+=3;
	}

	diag("big increment");
	j = 0;
	for ( i = 0 ; i < 5 ; i++, j++){
		pass("%d %d", i, j);
	}

	diag("big increment");
	i = 0;
	j = 0;
	for ( ; i < 5 ; i++, j++){
		pass("%d %d", i, j);
	}
*/
//	diag("big ininitialization and increment");
//	for (/*comment*/ i = 0 /*comment*/,/*comment*/ j = 0 /*comment*/; 
//			i </*comment*/ 5 ;
//		   	i++ /* comment*/ , f++ /*comment*/
//			){
//		pass("%d %d", i, j);
//	}

    done_testing();
}
