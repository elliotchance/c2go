#include <stdio.h>

int main()
{
    int value = 1;

    while (value <= 3)
    {
        printf("Value is %d\n", value);
        value++;
    }

    // continue
    value = 0;
    while (value < 3) {
	value++;
        if (value < 3) continue;
        printf("%d\n", value);
    }

    return 0;
}
