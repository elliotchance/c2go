#include <stdio.h>
#include "tests.h"

int main()
{
    plan(19);

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

    done_testing();
}
