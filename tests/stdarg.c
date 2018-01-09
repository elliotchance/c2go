#include <stdio.h>
#include <stdarg.h>
#include <assert.h>
#include "tests.h"

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void simple(const char* fmt, ...)
{
	char buffer[155];
	for (int i=0;i<155;i++){
		buffer[i] = 0;
	}
	is_streq(buffer, "");

    va_list args;
    va_start(args, fmt);
			
	char temp[100];
	for (int i=0;i<100;i++){
		temp[i] = 0;
	}

	int len = 4;
	for (int i=0;i<len;i++){
		char f = fmt[i];
        if (f == 'd') {
            int i = va_arg(args, int);
			sprintf(temp, "%d ",i);
			strcat(buffer , temp);
			is_streq(buffer, "3 ")
        } else if (f == 'c') {
            // note automatic conversion to integral type
            int c = va_arg(args, int);
			sprintf(temp, "%c ",c);
			strcat(buffer , temp);
			is_streq(buffer, "3 a ")
        } else if (f == 'f') {
            double d = va_arg(args, double);
			sprintf(temp, "%.3f ",d);
			strcat(buffer , temp);
        }
    }
 
    va_end(args);

	is_streq(buffer, "3 a 1.999 42.500 ")
}

void test_va_list()
{
	simple("dcff", 3, 'a', 1.999, 42.5); 
}

int main()
{
    plan(4);

    START_TEST(va_list)

    done_testing();
}
