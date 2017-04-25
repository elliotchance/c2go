#include <stdio.h>
#include <stdbool.h>

int main()
{
    bool trueBool = true;
    bool falseBool = false;

    if (trueBool)
    {
        printf("true bool\n");
    }

    if (!falseBool)
    {
        printf("reversed false bool\n");
    }
}
