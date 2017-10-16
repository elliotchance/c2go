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

int main()
{
    plan(8);

    struct programming variable;
    char *s = "Programming in Software Development.";

    variable.constant = 1.23;
    variable.pointer = s;

    is_eq(variable.constant, 1.23);
    is_streq(variable.pointer, "Programming in Software Development.");

    pass_by_val(variable);
    pass_by_ref(&variable);
/*
    |-DeclStmt 0x3b12420 <line:52:5, col:25>
    | `-VarDecl 0x3b123c0 <col:5, col:23> col:23 used s1 'struct mainStruct':'struct mainStruct'
    |-BinaryOperator 0x3b124d0 <line:53:5, col:19> 'float' '='
    | |-MemberExpr 0x3b12460 <col:5, col:8> 'float' lvalue .constant 0x3b0fe30
    | | `-DeclRefExpr 0x3b12438 <col:5> 'struct mainStruct':'struct mainStruct' lvalue Var 0x3b123c0 's1' 'struct mainStruct':'struct mainStruct'
    | `-ImplicitCastExpr 0x3b124b8 <col:19> 'float' <FloatingCast>
    |   `-FloatingLiteral 0x3b12498 <col:19> 'double' 4.200000e+01
*/
    struct mainStruct s1;
    s1.constant = 42.;
    is_eq(s1.constant, 42.);
/*
   |-DeclStmt 0x3b134f0 <line:56:2, col:17>
    | `-VarDecl 0x3b13490 <col:2, col:15> col:15 used s2 'secondStruct':'struct mainStruct'
    |-BinaryOperator 0x3b135a0 <line:57:2, col:16> 'float' '='
    | |-MemberExpr 0x3b13530 <col:2, col:5> 'float' lvalue .constant 0x3b0fe30
    | | `-DeclRefExpr 0x3b13508 <col:2> 'secondStruct':'struct mainStruct' lvalue Var 0x3b13490 's2' 'secondStruct':'struct mainStruct'
    | `-ImplicitCastExpr 0x3b13588 <col:16> 'float' <FloatingCast>
    |   `-FloatingLiteral 0x3b13568 <col:16> 'double' 4.200000e+01
*/
	secondStruct s2;
	s2.constant = 42.;
	is_eq(s2.constant, 42.);

    done_testing();
}
