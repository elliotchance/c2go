#include "tests.h"
#include <stdio.h>

void f_empty()
{
    return;
};

int main()
{
    plan(9);

    int a = 'a' == 65 ? 10 : 100;
    float b = 10 == 10 ? 1.0 : 2.0;
    char* c = 'x' == 5 ? "one" : "two";
    char d = a == 100 ? 'x' : 1;

    is_eq(a, 100);
    is_eq(b, 1);
    is_streq(c, "two");
    is_eq(d, 'x');

    is_false(0 ? 1 : 0);
    is_false(NULL ? 1 : 0);
    is_true('x' ? 1 : 0);

    a = a == 10 ? b == 1.0 ? 1 : 2 : 2;

    if (a == (a == 2 ? 5 : 10)) {
        fail(__func__);
    } else {
        pass(__func__);
    }

    diag("CStyleCast <ToVoid>");
    {
        double a, b;
        0 ? (void)(a) : (void)(b);
        (void)(a), (void)(b);
    }
    {
        double a;
        0 ? (void)(a) : f_empty();
        (void)(a);
    }
    {
        double b;
        0 ? f_empty() : (void)(b);
        (void)(b);
    }
    {
        ;
        0 ? f_empty() : f_empty();
    }
    pass("Ok - ToVoid");

    done_testing();
}
