#include <stdio.h>
#include <syslog.h>
#include "tests.h"

int main()
{
	plan(0);
	openlog("c2go",LOG_PID|LOG_NDELAY,LOG_USER);
	setlogmask(LOG_UPTO(LOG_NOTICE));
	int i = 1;
	double x = 2.71828;
	syslog(LOG_NOTICE,"hi there %d %f",i,x);
	closelog();
	done_testing();
}

