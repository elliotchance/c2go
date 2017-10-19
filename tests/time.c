#include <stdio.h>
#include <time.h>

#include "tests.h"

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_time()
{
    time_t now;
    time_t tloc;

    now = time(NULL);

    now = time(&tloc);
    is_true(now == tloc);

    pass("%s", "time");
}

void test_ctime()
{
    char* s = ctime(NULL);
    is_null(s);

    time_t now = 946684798;
    s = ctime(&now);
    printf("%s", s);

    pass("%s", "ctime");
}

int main()
{
    plan(4);

    START_TEST(time)
    START_TEST(ctime)

    done_testing();
}
