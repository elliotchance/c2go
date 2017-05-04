// Array examples

#include <stdio.h>

int main()
{
    int array[3], n, c;

    n = 3;

    array[0] = 5;
    array[1] = 9;
    array[2] = -13;

    printf("Array elements entered by you are:\n");

    for (c = 0; c < n; c++)
        printf("array[%d] = %d\n", c, array[c]);

    return 0;
}
