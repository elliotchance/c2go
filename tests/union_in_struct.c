// Tests for unions.

#include <stdio.h>
#include "tests.h"

typedef struct FuncDestructor FuncDestructor;
struct FuncDestructor {
	int i;
};
typedef struct FuncDef FuncDef;
struct FuncDef {
  int i;
  union {
    FuncDef *pHash;
    FuncDestructor *pDestructor;
  } u;
};

void union_inside_struct2()
{
	FuncDef f;
	FuncDestructor fd;
	fd.i = 100;
	f.u.pDestructor = &fd;

	FuncDestructor * p_fd = f.u.pDestructor;
	is_eq((*p_fd).i , 100);

	is_true(f.u.pHash       != NULL);
	is_true(f.u.pDestructor != NULL);
	int vHash = (*f.u.pHash).i;
	is_eq(vHash          , 100);
	is_eq((*f.u.pHash).i , 100);
}


union programming
{
    float constant;
    char *pointer;
};

union programming init_var()
{
    union programming variable;
    char *s = "Programming in Software Development.";

    variable.pointer = s;
    is_streq(variable.pointer, "Programming in Software Development.");

    variable.constant = 1.23;
    is_eq(variable.constant, 1.23);

    return variable;
}

void pass_by_ref(union programming *addr)
{
    char *s = "Show string member.";
    float v = 1.23+4.56;

    addr->constant += 4.56;
    is_eq(addr->constant, v);

    addr->pointer = s;
    is_streq(addr->pointer, "Show string member.");
}

void var_by_val(union programming value)
{
    value.constant++;

    is_eq(value.constant, 2.23);
}

int main()
{
    plan(16);
    
	union programming variable;

    variable = init_var();
    var_by_val(variable);
    pass_by_ref(&variable);
	
	union_inside_struct();
	union_inside_struct2();
    
	done_testing();
}
