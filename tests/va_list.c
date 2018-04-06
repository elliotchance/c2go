#include <stdio.h>
#include <stdarg.h>
#include <assert.h>
#include "tests.h"

void simple(const char* fmt, va_list args)
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

void test_va_list(const char *format, ...)
{
	va_list args;
	va_start(args, format);
	simple(format, args);
	va_end(args);
}

void test_va_list2(void (*f) (const char *format, va_list args), char *format, ...) {
	va_list args;
	va_start(args, format);
	f(format, args);
	va_end(args);
}

typedef void func_va (const char *format, va_list args);

void test_va_list3(func_va *f, char *format, ...) {
	va_list args;
	va_start(args, format);
	f(format, args);
	va_end(args);
}

int main()
{
	plan(12);

	test_va_list("dcff", 3, 'a', 1.999, 42.5);
	test_va_list2(simple, "dcff", 3, 'a', 1.999, 42.5);
	test_va_list3(simple, "dcff", 3, 'a', 1.999, 42.5);

	done_testing();
}
