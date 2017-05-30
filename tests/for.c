#include <stdio.h>
#include "tests.h"

int main()
{
    int i = 0;

    // Missing init
    for (; i < 10; i++)
        printf("%d\n", i);

    // CompountStmt
    for (i = 0; i < 10; i++) {
        printf("%d\n", i);
    }

    // Not CompoundStmt
    for (i = 0; i < 10; i++)
        printf("%d\n", i);

    // Infinite loop
    int j = 0;
    for (;;) {
        printf("infinite loop\n");
        j++;
        if (j > 10)
            break;
    }

    // continue
    i = 0; j = 0;
    for(;;) {
        printf("%d %d\n", i, j);
        i++;
        if (i < 3) continue;
        j++;
        if (j > 3) break;
    }
}
