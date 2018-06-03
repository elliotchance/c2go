#include <string.h>
#include "tests.h"

typedef struct mem {
    int a;
    int b;
} mem;

typedef struct mem2 {
    int a[2];
} mem2;
typedef int altint;

void setptr(int *arr, int val) {
    arr[0] = val;
}
void setarr(int arr[], int val) {
    arr[0] = val;
}

int main()
{
    plan(80);

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
        diag("strncmp");
        {
            char* a = "ab";
            char* b = "ab";
            is_true(strncmp(a,b,10) == 0);
        }
        {
            char* a = "bb";
            char* b = "ab";
            is_true(strncmp(a,b,10) > 0);
        }
        {
            char* a = "ab";
            char* b = "bb";
            is_true(strncmp(a,b,10) < 0);
        }
        {
            char* a = "aba";
            char* b = "ab";
            is_true(strncmp(a,b,10) > 0);
        }
        {
            char* a = "ab";
            char* b = "aba";
            is_true(strncmp(a,b,10) < 0);
        }
        {
            char* a = "aba";
            char* b = "abc";
            is_true(strncmp(a,b,2) == 0);
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
    {
        diag("memset");
        char dest1[40];
        char *dest2;
        char *dest3;
        dest2 = (char*) memset(dest1, 'a', 4);
        dest1[4] = '\0';
        is_streq(dest1, "aaaa");
        is_streq(dest2, "aaaa");
        dest3 = (char*) memset(&dest2[1], 'b', 2);
        is_streq(dest1, "abba");
        is_streq(dest2, "abba");
        is_streq(dest3, "bba");
    }
    {
        diag("memcpy");
        char *src = "aaaabb";
        char dest1[40];
        char *dest2;
        char *dest3;
        dest2 = (char*) memcpy(dest1, src, 4);
        dest1[4] = '\0';
        is_streq(dest1, "aaaa");
        is_streq(dest2, "aaaa");
        dest3 = (char*) memcpy(&dest2[1], &src[4], 2);
        is_streq(dest1, "abba");
        is_streq(dest2, "abba");
        is_streq(dest3, "bba");
    }
    {
        diag("memmove");
        char *src = "aaaabb";
        char dest1[40];
        char *dest2;
        char *dest3;
        dest2 = (char*) memmove(dest1, src, 4);
        dest1[4] = '\0';
        is_streq(dest1, "aaaa");
        is_streq(dest2, "aaaa");
        dest3 = (char*) memmove(&dest2[1], &src[4], 2);
        is_streq(dest1, "abba");
        is_streq(dest2, "abba");
        is_streq(dest3, "bba");
    }
    {
        diag("memset & memcpy of struct / array of struct");
        int dest3 = 4;
        int dest4 = 0xf;
        memcpy(&dest3, &dest4, sizeof(int));
        is_eq(dest3, 0xf);
        memset(&dest4, 0, sizeof(int));
        is_eq(dest4, 0);
        mem dest5 = {
            .a = 2,
            .b = 3.0
        };
        memset(&dest5, 0, sizeof(mem));
        is_eq(dest5.a, 0);
        is_eq(dest5.b, 0.0);
        mem dest6[] = {
            {
               .a = 2,
               .b = 3.0
            },
            {
               .a = 4,
               .b = 5.0
            }
        };
        mem dest7[2];
        memcpy(dest7, dest6, sizeof(mem)*2);
        memset(dest6, 0, sizeof(mem)*2);
        is_eq(dest6[0].a, 0);
        is_eq(dest6[0].b, 0.0);
        is_eq(dest6[1].a, 0);
        is_eq(dest6[1].b, 0.0);
        is_eq(dest7[0].a, 2);
        is_eq(dest7[0].b, 3.0);
        is_eq(dest7[1].a, 4);
        is_eq(dest7[1].b, 5.0);
        memset(&dest7[1], 0, sizeof(mem));
        is_eq(dest7[0].a, 2);
        is_eq(dest7[0].b, 3.0);
        is_eq(dest7[1].a, 0);
        is_eq(dest7[1].b, 0.0);
        mem2 dest8;
        dest8.a[0] = 42;
        memset(dest8.a, 0, sizeof(int)*2);
        is_eq(dest8.a[0], 0);
        is_eq(dest8.a[1], 0);
        dest8.a[0] = 42;
        mem2 dest9;
        altint *test = (altint *) dest9.a;
        memcpy(dest9.a, dest8.a, sizeof(int)*2);
        is_eq(dest9.a[0], 42);
        is_eq(test[0], 42);
        setarr(dest9.a, 1);
        is_eq(dest9.a[0], 1);
        setptr(dest9.a, 2);
        is_eq(dest9.a[0], 2);
    }
    {
        diag("memcmp");
        {
            char* a = "ab\0c";
            char* b = "ab\0c";
            is_true(memcmp(a,b,4) == 0);
        }
        {
            char* a = "ab\0a";
            char* b = "ab\0c";
            is_true(memcmp(a,b,4) < 0);
        }
        {
            char* a = "ab\0c";
            char* b = "ab\0a";
            is_true(memcmp(a,b,4) > 0);
        }
        {
            char* a = "ab\0c";
            char* b = "ab\0a";
            is_true(memcmp(a,b,3) == 0);
        }
    }
    {
        diag("strstr");
        {
            char* a = "needle in a haystack";
            char* b = "haystack";
            char* res = strstr(a,b);
            is_streq(res, "haystack");
        }
        {
            char* a = "needle in a haystack";
            char* b = "wrong";
            char* res = strstr(a,b);
            is_null(res);
        }
    }
    {
        diag("strcasestr");
        {
            char* a = "needle in a haystack";
            char* b = "HayStack";
            char* res = strcasestr(a,b);
            is_streq(res, "haystack");
        }
        {
            char* a = "needle in a haystack";
            char* b = "wrong";
            char* res = strcasestr(a,b);
            is_null(res);
        }
    }

    done_testing();
}
