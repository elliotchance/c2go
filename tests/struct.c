// Tests for structures.

#include <stdio.h>
#include "tests.h"

struct programming
{
    float constant;
    char *pointer;
};

void pass_by_ref(struct programming *addr)
{
    char *s = "Show string member.";
    float v = 1.23+4.56;

    addr->constant += 4.56;
    addr->pointer = s;

    is_eq(addr->constant, v);
    is_streq(addr->pointer, "Show string member.");
}

void pass_by_val(struct programming value)
{
    value.constant++;

    is_eq(value.constant, 2.23);
    is_streq(value.pointer, "Programming in Software Development.");
}

typedef struct mainStruct{
	double constant;
} secondStruct;

typedef struct {
	double t;
} ts_c;

typedef struct ff {
	int v1,v2;
} tt1, tt2;

int main()
{
    plan(12);

    struct programming variable;
    char *s = "Programming in Software Development.";

    variable.constant = 1.23;
    variable.pointer = s;

    is_eq(variable.constant, 1.23);
    is_streq(variable.pointer, "Programming in Software Development.");

    pass_by_val(variable);
    pass_by_ref(&variable);
    
	struct mainStruct s1;
    s1.constant = 42.;
    is_eq(s1.constant, 42.);
	
	secondStruct s2;
	s2.constant = 42.;
	is_eq(s2.constant, 42.);
	
	ts_c c;
	c.t = 42.;
	is_eq(c.t , 42.);

	tt1 t1;
	t1.v1 = 42.;
	is_eq(t1.v1,42.)

	tt2 t2;
	t2.v1 = 42.;
	is_eq(t2.v1,42.)

	struct ff tf2;
	tf2.v2 = t1.v1;
	is_eq(tf2.v2,t1.v1);

    done_testing();
}
