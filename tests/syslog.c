#include <stdio.h>
#include <syslog.h>
#include "tests.h"

int main()
{
	plan(0);
	openlog("c2go",LOG_PID|LOG_NDELAY,LOG_USER);
	setlogmask(LOG_UPTO(LOG_NOTICE));
	syslog(LOG_NOTICE,"hi there");
	closelog();
	done_testing();
}

