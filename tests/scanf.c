// scanf() needs to be in it's own file because it takes input from stdin.

#include <stdio.h>
#include "tests.h"

int main()
{
    plan(1);

    int i;
    scanf("%d", &i);
    is_eq(i, 7);

    done_testing();
}
