#include <string.h>
#include "tests.h"

int main()
{
    plan(29);

    diag("TODO: __builtin_object_size")
    // https://github.com/elliotchance/c2go/issues/359

    {
        diag("strcpy")
        char *src = "foo bar\0baz";
        char dest1[40];
        char *dest2;
        dest2 = strcpy(dest1, src);
        is_streq(dest1, "foo bar");
        is_streq(dest2, "foo bar");
    }

    diag("strncpy")

    // If the end of the source C string (which is signaled by a null-character)
    // is found before num characters have been copied, destination is padded
    // with zeros until a total of num characters have been written to it.
    {
        char *src = "foo bar\0baz";
        char dest1[40];
        char *dest2;

        dest1[7] = 'a';
        dest1[15] = 'b';
        dest1[25] = 'c';
        dest2 = strncpy(dest1, src, 20);

        is_eq(dest1[0], 'f');
        is_eq(dest1[7], 0);
        is_eq(dest1[15], 0);
        is_eq(dest1[25], 'c');

        is_eq(dest2[0], 'f');
        is_eq(dest2[7], 0);
        is_eq(dest2[15], 0);
        is_eq(dest2[25], 'c');

        is_streq(dest1, "foo bar");
        is_streq(dest2, "foo bar");
    }

    // No null-character is implicitly appended at the end of destination if
    // source is longer than num. Thus, in this case, destination shall not be
    // considered a null terminated C string (reading it as such would
    // overflow).
    {
        char *src = "foo bar\0baz";
        char dest1[40];
        char *dest2;

        dest1[7] = 'a';
        dest1[15] = 'b';
        dest1[25] = 'c';
        dest2 = strncpy(dest1, src, 5);

        is_eq(dest1[0], 'f');
        is_eq(dest1[7], 'a');
        is_eq(dest1[15], 'b');
        is_eq(dest1[25], 'c');

        is_eq(dest2[0], 'f');
        is_eq(dest2[7], 'a');
        is_eq(dest2[15], 'b');
        is_eq(dest2[25], 'c');
    }

    {
        diag("strlen")
        is_eq(strlen(""), 0);
        is_eq(strlen("a"), 1);
        is_eq(strlen("foo"), 3);
        // NULL causes a seg fault.
        // is_eq(strlen(NULL), 0);
        is_eq(strlen("fo\0o"), 2);
    }
	{
		diag("strcat")
		char str[80];
		strcpy (str,"these ");
		strcat (str,"strings ");
		strcat (str,"are ");
		strcat (str,"concatenated.");
		is_streq(str,"these strings are concatenated.");
	}
	{
		diag("strcmp");
		{
			char* a = "ab";
			char* b = "ab";
			is_true(strcmp(a,b) == 0);
		}
		{
			char* a = "bb";
			char* b = "ab";
			is_true(strcmp(a,b) > 0);
		}
		{
			char* a = "ab";
			char* b = "bb";
			is_true(strcmp(a,b) < 0);
		}
	}
	{
		diag("strchr");
		char str[] = "This is a sample string";
		char * pch;
		int amount = 0;
		pch=strchr(str,'s');
		while (pch!=NULL)
		{
			pch=strchr(pch+1,'s');
			amount ++;
		}
		is_eq(amount,  4 );
	}

    done_testing();
}
