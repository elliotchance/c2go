#include <string.h>
#include "tests.h"

int main()
{
    plan(5);

    diag("strlen")
    is_eq(strlen(""), 0);
    is_eq(strlen("a"), 1);
    is_eq(strlen("foo"), 3);
    is_eq(strlen(NULL), 0);
    is_eq(strlen("fo\0o"), 2);

    done_testing();
}
