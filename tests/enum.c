#include <stdio.h>
#include "tests.h"

enum week{Mon, Tue, Wed};

int main()
{
	plan(1);
    
	enum week day;
	day = Wed;
	is_eq(day ,2)

	done_testing();
}
