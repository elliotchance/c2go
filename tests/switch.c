// This file tests the various forms of switch statement.
//
// We must be extra sensitive to the fact that switch fallthrough is handled
// differently in C and Go. Break statements are removed and fallthrough
// statements are added when necessary.
//
// It is worth mentioning that a SwitchStmt has a CompoundStmt item that
// contains all of the cases. However, if the individual case are not enclosed
// in a scope each of the case statements and their childen are part of the same
// CompoundStmt. For example, the first switch statement below contains a
// CompoundStmt with 12 children.

#include <stdio.h>
#include <stdbool.h>
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

void fallthrough_several_midway_default()
{
    for(int i=0; i<=3; i++) {
        int j = -1;
        int expected = -1;
        if(i==0)
            expected = 10;
        if(i==1)
            expected = 21;
        if(i==2)
            expected = 22;
        if(i==3)
            expected = 13;
        if(i==4)
            expected = 666;
        switch(i) {
        case 4:
            fail("code should not reach here");
        case 0:
        default:
            j = i+10;
            break;
        case 1:
        case 2:
            j = i+20;
            break;
        }
        is_eq(j, expected);
    }
}

void goto_label(bool use_goto)
{
    for (;;) {
        switch (0)
        {
        case 3:
            continue;
        case 0:
            if (use_goto) {
                for (;;)
                    break;
                goto LABEL;
                fail("code should not reach here");
            } else if (false) {
                goto LABELX;
                goto LABELY;
                fail("code should not reach here");
            }
            /* other comment */
            // some comment
            /* fallthrough */
        LABELY:
        case 4:
        LABEL:
            printf("x");
        case 1:
            pass(__func__);
            break;
        case 2:
            ;
        LABELX:
        default:
            fail("code should not reach here");
            break;
        }
        break;
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

void scoped_fallthrough_several_midway_default()
{
    for(int i=0; i<=3; i++) {
        int j = -1;
        int expected = -1;
        if(i==0)
            expected = 10;
        if(i==1)
            expected = 21;
        if(i==2)
            expected = 22;
        if(i==3)
            expected = 13;
        if(i==4)
            expected = 666;
        switch(i) {
        case 4:
        {
            fail("code should not reach here");
        }
        case 0:
        {}
        default:
        {
            j = i+10;
            break;
        }
        case 1:
        {}
        case 2:
        {
            j = i+20;
        }
        }
        is_eq(j, expected);
    }
}

void scoped_goto_label(bool use_goto)
{
    for (;;) {
        switch (0)
        {
        case 3:
            {
                continue;
            }
        case 0:
            {
                if (use_goto) {
                    for (;;) {
                        break;
                    }
                    goto LABEL;
                    fail("code should not reach here");
                } else if (false) {
                    goto LABELX;
                    goto LABELY;
                    fail("code should not reach here");
                }
                /* other comment */
                // some comment
                /* fallthrough */
            }
        LABELY: {}
        case 4: {}
        LABEL:
            {
                printf("x");
            }
        case 1:
            {
                pass(__func__);
                break;
            }
        case 2:
            {
                int x = 0;
                printf("%d", x);
                break;
            }
        LABELX: {}
        default:
            {
                fail("code should not reach here");
                break;
            }
        }
        break;
    }
}

typedef struct I67 I67;
struct I67{
	int x,y;
};

void run( I67 * i ,int pos)
{
	switch (pos) {
		case 0:
			(*i).x += 1;
			(*i).y += 1;
			break;
		case 1:
			(*i).x -= 1;
			(*i).y -= 1;
			break;
	}
}

void run_with_block( I67 * i ,int pos)
{
	switch (pos) {
		case 0:
			{
			(*i).x += 1;
			(*i).y += 1;
			break;
			}
		case 1:
			{
			(*i).x -= 1;
			(*i).y -= 1;
			}
			break;
		case 2:
			(*i).x *= 1;
			(*i).y *= 1;
			break;
		default:
			(*i).x *= 10;
			(*i).y *= 10;
	}
}

void switch_issue67()
{
	I67 i;
	i.x = 0;
	i.y = 0;
	run(&i, 0);
	is_eq(i.x, 1);
	is_eq(i.y, 1);
	run(&i, 1);
	is_eq(i.x, 0);
	is_eq(i.y, 0);
	run_with_block(&i,0);
	is_eq(i.x, 1);
	is_eq(i.y, 1);
	run_with_block(&i, 1);
	is_eq(i.x, 0);
	is_eq(i.y, 0);
}

void empty_switch()
{
	int pos = 0;
	switch (pos){
	}
	is_eq(pos,0);
}

void default_only_switch()
{
	int pos = 0;
	switch (pos){
		case -1: // empty case
		case -1-4: // empty case
		case (-1-4-4): // empty case
		case (-3): // empty case
		case -2: // empty case
		default:
			pos++;
	}
	is_eq(pos,1);
}

void switch_without_input()
{
	int pos = 0;
	switch (0){
		default:
			pos++;
	}
	is_eq(pos,1);
}

int main()
{
    plan(37);

    match_a_single_case();
    fallthrough_to_next_case();
    match_no_cases();
    match_default();
    fallthrough_several_cases_including_default();
    fallthrough_several_midway_default();
    goto_label(false);
    goto_label(true);

    // For each of the tests above there will be identical cases that use scopes
    // for the case statements.
    scoped_match_a_single_case();
    scoped_fallthrough_to_next_case();
    scoped_match_no_cases();
    scoped_match_default();
    scoped_fallthrough_several_cases_including_default();
    scoped_fallthrough_several_midway_default();
    scoped_goto_label(false);
    scoped_goto_label(true);

	switch_issue67();
	empty_switch();
	default_only_switch();
	switch_without_input();

    done_testing();
}
