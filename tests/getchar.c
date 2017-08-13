// getchar() needs to be in it's own file because it takes input from stdin.

#include <stdio.h>
#include "tests.h"

int main()
{
    plan(1);

    int c;
    c = getchar();
    is_eq(c, '7');

    done_testing();
}
