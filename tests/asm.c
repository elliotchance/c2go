// The functions in this file were adapted from:
// http://read.cs.ucla.edu/sqlite-anvil/trunk/sqlite-3.6.0/src/hwtime.h
//
// For more information see the original issue:
// https://github.com/elliotchance/c2go/issues/228

#include "tests.h"

__inline__ unsigned long sqlite3Hwtime1(void){
    unsigned int lo, hi;
    __asm__ __volatile__ ("rdtsc" : "=a" (lo), "=d" (hi));
    return (unsigned long)hi << 32 | lo;
}

// This has been disabled because clang cannot understand the asm{} syntax.
//__inline sqlite_uint64 __cdecl sqlite3Hwtime2(void){
//    __asm {
//        rdtsc
//        ret       ; return value at EDX:EAX
//    }
//}

__inline__ unsigned long sqlite3Hwtime3(void){
    unsigned long val;
    __asm__ __volatile__ ("rdtsc" : "=A" (val));
    return val;
}

int main() {
    // There are no actual tests in this file because Go does not support inline
    // assembly. We will have to revisit this in the future.
    plan(0);

    done_testing();
}
