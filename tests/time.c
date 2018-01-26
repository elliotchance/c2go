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

void test_mktime()
{
  time_t rawtime;
  struct tm * timeinfo;
  int year, month ,day;

  year  = 2000;
  month = 5;
  day   = 20;

  /* get current timeinfo and modify it to the user's choice */
  time ( &rawtime );
  timeinfo = localtime ( &rawtime );
  timeinfo->tm_year = year - 1900;
  timeinfo->tm_mon = month - 1;
  timeinfo->tm_mday = day;

  /* call mktime: timeinfo->tm_wday will be set */
  mktime ( timeinfo );

  is_eq(timeinfo->tm_wday,6);
}

int main()
{
    plan(6);

    START_TEST(time);
    START_TEST(ctime);
	START_TEST(mktime);

    done_testing();
}
