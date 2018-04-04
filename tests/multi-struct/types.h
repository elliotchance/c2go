#ifndef HEADER
#define HEADER

struct real_type;
typedef struct real_type type;

typedef struct mystery {
  unsigned int answer;
} mystery;

// Forward-declared prototypes that are defined in one of our other C files.
void say_type(type *t);   // types.c
type* to_type(mystery m); // types.c

#endif /* HEADER */
