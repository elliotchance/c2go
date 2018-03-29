#include <stdio.h>
#include <stdlib.h>
#include "types.h"

struct real_type {
    unsigned int real_answer;
};

void say_type(type *t) {
    printf("%u\n", t->real_answer);
}

type* to_type(mystery m) {
    type *ret = (type*) malloc(sizeof(type));
    ret->real_answer = m.answer;
    return ret;
}
