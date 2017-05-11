// This file tests the various forms of switch statement.
//
// We must be extra sensitive to the fact that switch fallthrough is handled
// differently in C and Go. Break statements are removed and fallthrough
// statements are added when nessesary.
//
// It is worth mentioning that a SwitchStmt has a CompoundStmt item that
// contains all of the cases. However, if the individual case are not enclosed
// in a scope each of the case statements and their childen are part of the same
// CompoundStmt. For example, the first switch statement below contains a
// CompoundStmt with 12 children.

#include <stdio.h>

void match_a_single_case()
{
    switch (1)
    {
    case 0:
        printf("a 00\n");
        printf("a 01\n");
        printf("a 02\n");
        break;
    case 1:
        printf("a 10\n");
        printf("a 11\n");
        printf("a 12\n");
        break;
    case 2:
        printf("a 20\n");
        printf("a 21\n");
        printf("a 22\n");
        break;
    default:
        printf("a default 0\n");
        printf("a default 1\n");
        printf("a default 2\n");
        break;
    }
}

void fallthrough_to_next_case()
{
    switch (1)
    {
    case 0:
        printf("b 0\n");
        break;
    case 1:
        printf("b 1\n");
    case 2:
        printf("b 2\n");
        break;
    default:
        printf("b default\n");
        break;
    }
}

void match_no_cases()
{
    switch (1)
    {
    case 5:
        printf("c 5\n");
        break;
    case 2:
        printf("c 2\n");
        break;
    }
}

void match_default()
{
    switch (1)
    {
    case 5:
        printf("d 5\n");
        break;
    case 2:
        printf("d 2\n");
        break;
    default:
        printf("d default\n");
        break;
    }
}

void fallthrough_several_cases_including_default()
{
    switch (1)
    {
    case 0:
        printf("e 0\n");
    case 1:
        printf("e 1\n");
    case 2:
        printf("e 2\n");
    default:
        printf("e default\n");
    }
}

void scoped_match_a_single_case()
{
    switch (1)
    {
    case 0:
    {
        printf("a 0\n");
        break;
    }
    case 1:
    {
        printf("a 1\n");
        break;
    }
    case 2:
    {
        printf("a 2\n");
        break;
    }
    default:
    {
        printf("a default\n");
        break;
    }
    }
}

void scoped_fallthrough_to_next_case()
{
    switch (1)
    {
    case 0:
    {
        printf("b 0\n");
        break;
    }
    case 1:
    {
        printf("b 1\n");
    }
    case 2:
    {
        printf("b 2\n");
        break;
    }
    default:
    {
        printf("b default\n");
        break;
    }
    }
}

void scoped_match_no_cases()
{
    switch (1)
    {
    case 5:
    {
        printf("c 5\n");
        break;
    }
    case 2:
    {
        printf("c 2\n");
        break;
    }
    }
}

void scoped_match_default()
{
    switch (1)
    {
    case 5:
    {
        printf("d 5\n");
        break;
    }
    case 2:
    {
        printf("d 2\n");
        break;
    }
    default:
    {
        printf("d default\n");
        break;
    }
    }
}

void scoped_fallthrough_several_cases_including_default()
{
    switch (1)
    {
    case 0:
    {
        printf("e 0\n");
    }
    case 1:
    {
        printf("e 1\n");
    }
    case 2:
    {
        printf("e 2\n");
    }
    default:
    {
        printf("e default\n");
    }
    }
}

int main()
{
    match_a_single_case();
    fallthrough_to_next_case();
    match_no_cases();
    match_default();
    fallthrough_several_cases_including_default();

    // For each of the tests above there will be identical cases that use scopes
    // for the case statements.
    scoped_match_a_single_case();
    scoped_fallthrough_to_next_case();
    scoped_match_no_cases();
    scoped_match_default();
    scoped_fallthrough_several_cases_including_default();
}
