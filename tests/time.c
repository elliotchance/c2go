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

void test_localtime()
{
	struct tm * timeinfo;
	time_t      rawtime = 80000;
	timeinfo = localtime ( &rawtime );
	is_eq( timeinfo-> tm_sec	,  20 );
	is_eq( timeinfo-> tm_min	,  13 );
	is_eq( timeinfo-> tm_hour	,  1  );
	is_eq( timeinfo-> tm_mday	,  2  );
	is_eq( timeinfo-> tm_mon	,  0  );
	is_eq( timeinfo-> tm_year	,  70 );
	is_eq( timeinfo-> tm_wday	,  5  );
	is_eq( timeinfo-> tm_yday	,  1  );
	is_eq( timeinfo-> tm_isdst	,  0  );
}

void test_mktime()
{
	struct tm  timeinfo;
	
	timeinfo.tm_year = 2000  - 1900;
	timeinfo.tm_mon  = 5     - 1   ;
	timeinfo.tm_mday = 20          ;
	timeinfo.tm_sec  = 0           ;
	timeinfo.tm_min  = 0           ;
	timeinfo.tm_hour = 0           ;
	
	time_t t = mktime ( &timeinfo );
	
	is_eq(timeinfo.tm_wday  , 6           );
	is_eq(timeinfo.tm_year  , 100         );
	is_eq(timeinfo.tm_mon   , 4           );
	is_eq(timeinfo.tm_mday  , 20          );
	is_eq(t                 , 958766400   );
}

int main()
{
	plan(19);
	
	START_TEST(time      );
	START_TEST(ctime     );
	START_TEST(localtime );
	START_TEST(mktime    );
	
	done_testing();
}
