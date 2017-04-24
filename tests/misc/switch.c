#include <stdio.h>

int main()
{
    int num = 1;

    // Match the single case.
    switch (num)
    {
    case 0:
        printf("a 0\n");
        break;
    case 1:
        printf("a 1\n");
        break;
    case 2:
        printf("a 2\n");
        break;
    default:
        printf("a default\n");
        break;
    }

    // Fallthrough to next case.
    switch (num)
    {
    case 0:
        printf("b 0\n");
        break;
    case 1:
        printf("b 1\n");
    case 2:
        printf("b 2\n");
        break;
    default:
        printf("b default\n");
        break;
    }
}
