#include <stdlib.h>
#include "types.h"

int main() {
    mystery m;
    m.answer = 42;
    type* t = to_type(m);
    say_type(t);
    free(t);
    return 0;
}
