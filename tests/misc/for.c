#include <stdio.h>

int main()
{
    int i;
    for (i = 0; i < 10; i++)
        printf("%d\n", i);

    int j = 0;
    for (;;) {
        printf("infinite loop\n");
        j++;
        if (j > 10)
            break;
    }
}
