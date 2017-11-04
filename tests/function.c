#include <stdio.h>
#include "tests.h"

void my_function();

void function(){
}

int main()
{
    plan(4);

    pass("%s", "Main function.");

    my_function();

    pass("%s", "Back in function main.");

	// pointer on function
	void * a;
	a = function;
	is_not_null(a);	

    done_testing();
}

void my_function()
{
    pass("%s", "Welcome to my function. Feel at home.");
}
