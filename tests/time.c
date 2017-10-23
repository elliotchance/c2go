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
    is_not_eq(now, 0);

    now = time(&tloc);
    is_not_eq(now, 0);
    is_eq(now, tloc);
}

void test_ctime()
{
    char* s;

    // 1999-12-31 11:59:58
    time_t now = 946670398;
    s = ctime(&now);
    is_not_null(s);

    // Hours/minutes will vary based on local time. Ignore them.
    s[11] = 'H';
    s[12] = 'H';
    s[14] = 'm';
    s[15] = 'm';
    is_streq(s, "Fri Dec 31 HH:mm:58 1999\n");
}

int main()
{
    plan(5);

    START_TEST(time)
    START_TEST(ctime)

    done_testing();
}
