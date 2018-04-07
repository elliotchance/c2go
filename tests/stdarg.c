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

void simple2(const char* fmt, va_list args)
{
	char buffer[155];
	for (int i=0;i<155;i++){
		buffer[i] = 0;
	}
	is_streq(buffer, "");

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

	is_streq(buffer, "3 a 1.999 42.500 ")
}

void test_va_list()
{
	simple("dcff", 3, 'a', 1.999, 42.5); 
}

void test_va_list2(const char *format, ...)
{
	va_list args;
	va_start(args, format);
	simple2(format, args);
	va_end(args);
}

void test_va_list3(void (*f) (const char *format, va_list args), char *format, ...) {
	va_list args;
	va_start(args, format);
	f(format, args);
	va_end(args);
}

typedef void func_va (const char *format, va_list args);

void test_va_list4(func_va *f, char *format, ...) {
	va_list args;
	va_start(args, format);
	f(format, args);
	va_end(args);
}

int main()
{
    plan(16);

    START_TEST(va_list)

    test_va_list2("dcff", 3, 'a', 1.999, 42.5);
    test_va_list3(simple2, "dcff", 3, 'a', 1.999, 42.5);
    test_va_list4(simple2, "dcff", 3, 'a', 1.999, 42.5);

    done_testing();
}

