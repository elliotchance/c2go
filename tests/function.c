#include <stdio.h>
#include "tests.h"

void my_function();

int main()
{
    plan(3);

    pass("%s", "Main function.");

    my_function();

    pass("%s", "Back in function main.");

    done_testing();
}

void my_function()
{
    pass("%s", "Welcome to my function. Feel at home.");
}
