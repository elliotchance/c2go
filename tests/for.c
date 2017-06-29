#include <stdio.h>
#include "tests.h"

int main()
{
    plan(27);

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

	diag("big increment");
	i = 0;
	j = 0;
	for (;i < 2;i++,j++){
		pass("%d %d", i, j);
	}

	diag("bif condition");
	i = 0;
	j = 0;
	for(;i++,j<2;){
		pass("%d %d", i, j);
		j++;
	}

    done_testing();
}
