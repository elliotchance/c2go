#include <string.h>
#include "tests.h"

// https://gcc.gnu.org/onlinedocs/gcc/Object-Size-Checking.html
struct V { char buf1[10]; int b; char buf2[10]; };

int main()
{
    plan(6);

    diag("TODO: __builtin_object_size")
    // https://github.com/elliotchance/c2go/issues/359

    diag("strcpy")
    char *src = "foo bar\0baz";
    char dest1[40];
    char *dest2;
    dest2 = strcpy(dest1, src);
    is_streq(dest1, "foo bar");
    is_streq(dest2, "foo bar");

    diag("strlen")
    is_eq(strlen(""), 0);
    is_eq(strlen("a"), 1);
    is_eq(strlen("foo"), 3);
    // NULL causes a seg fault.
    // is_eq(strlen(NULL), 0);
    is_eq(strlen("fo\0o"), 2);

    done_testing();
}
