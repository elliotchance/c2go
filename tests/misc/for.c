#include <stdio.h>

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
}
