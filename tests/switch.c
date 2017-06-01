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
#include "tests.h"

void match_a_single_case()
{
    switch (1)
    {
    case 0:
        fail("code should not reach here");
        break;
    case 1:
        pass(__func__);
        break;
    case 2:
        fail("code should not reach here");
        break;
    default:
        fail("code should not reach here");
        break;
    }
}

void fallthrough_to_next_case()
{
    switch (1)
    {
    case 0:
        fail("code should not reach here");
        break;
    case 1:
        pass(__func__);
    case 2:
        pass(__func__);
        break;
    default:
        fail("code should not reach here");
        break;
    }
}

void match_no_cases()
{
    switch (1)
    {
    case 5:
        fail("code should not reach here");
        break;
    case 2:
        fail("code should not reach here");
        break;
    }
}

void match_default()
{
    switch (1)
    {
    case 5:
        fail("code should not reach here");
        break;
    case 2:
        fail("code should not reach here");
        break;
    default:
        pass(__func__);
        break;
    }
}

void fallthrough_several_cases_including_default()
{
    switch (1)
    {
    case 0:
        fail("code should not reach here");
    case 1:
        pass(__func__);
    case 2:
        pass(__func__);
    default:
        pass(__func__);
    }
}

void scoped_match_a_single_case()
{
    switch (1)
    {
    case 0:
    {
        fail("code should not reach here");
        break;
    }
    case 1:
    {
        pass(__func__);
        break;
    }
    case 2:
    {
        fail("code should not reach here");
        break;
    }
    default:
    {
        fail("code should not reach here");
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
        fail("code should not reach here");
        break;
    }
    case 1:
    {
        pass(__func__);
    }
    case 2:
    {
        pass(__func__);
        break;
    }
    default:
    {
        fail("code should not reach here");
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
        fail("code should not reach here");
        break;
    }
    case 2:
    {
        fail("code should not reach here");
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
        fail("code should not reach here");
        break;
    }
    case 2:
    {
        fail("code should not reach here");
        break;
    }
    default:
    {
        pass(__func__);
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
        fail("code should not reach here");
    }
    case 1:
    {
        pass(__func__);
    }
    case 2:
    {
        pass(__func__);
    }
    default:
    {
        pass(__func__);
    }
    }
}

int main()
{
    plan(14);

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

    done_testing();
}
