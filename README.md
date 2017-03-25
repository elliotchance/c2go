[![Build Status](https://travis-ci.org/elliotchance/c2go.svg?branch=master)](https://travis-ci.org/elliotchance/c2go)

A tool for converting C to Go.

# Why?

I created this project as a proof of concept. It is written in python and uses
the clang AST to do all of the hard work.

# How?

Let's use
[prime.c](https://github.com/elliotchance/c2go/blob/master/tests/prime.c):

```c
#include <stdio.h>
 
int main()
{
   int n, c;
 
   printf("Enter a number\n");
   scanf("%d", &n);
 
   if ( n == 2 )
      printf("Prime number.\n");
   else
   {
       for ( c = 2 ; c <= n - 1 ; c++ )
       {
           if ( n % c == 0 )
              break;
       }
       if ( c != n )
          printf("Not prime.\n");
       else
          printf("Prime number.\n");
   }
   return 0;
}
```

```bash
python c2go.py tests/prime.c
```

```go
package main

import (
    "fmt"
)

// ... lots of system types in Go removed for brevity.

func main() {
    var n int
    var c int
    fmt.Printf("Enter a number\n")
    fmt.Scanf("%d", &n)
    if n == 2 {
        fmt.Printf("Prime number.\n")
    } else {
        for c = 2; c <= n - 1; c += 1 {
            if n % c == 0 {
                break
            }
        }
        if c != n {
            fmt.Printf("Not prime.\n")
        } else {
            fmt.Printf("Prime number.\n")
        }
    }
    return
}
```

This is the process:

1. The C code is preprocessed with clang. This generates a larger file (`pp.c`),
but removes all the platform specific directives and macros.

2. The new file is parsed clang AST and dumps it in a colourful text format that
[looks like this](http://ehsanakhgari.org/wp-content/uploads/2015/12/Screen-Shot-2015-12-03-at-5.02.38-PM.png).
Apart from just parsing the C and dumping an AST, the AST contains all of the
resolved information that a compiler would need. This means that the code must
compile successfully under clang for the AST to also be usable.

3. Since we have all the types in the AST it's just a matter of traversing the
tree is a semi-intelligent way and producing Go. Easy, right!?

# Testing

Testing is done with a set of integrations tests in the form of complete C
programs that can be found in the
[tests](https://github.com/elliotchance/c2go/tree/master/tests) directory.

For each of those files:

1. Clang compiles the C to a binary as normal.
2. c2go converts the C file to Go.
3. The Go is built to produce another binary.
4. Both binaries are executed and the output is compared. All C files will
contain some output so the results can be verified.

The test suite is run with
[run-tests.sh](https://github.com/elliotchance/c2go/blob/master/run-tests.sh).

# Contributing

As I said it is still very early days (sorry for all the hacky Python). And
eventually the build chain can be converted to pure Go since we don't need any
clang APIs.

Contributing is done with pull requests. If you're looking for where to start I
can suggest
[finding a simple C program](http://www.programmingsimplified.com/c-program-examples)
(like the other examples) that does not successful translate into Go and fixing
up the Python so that it does.
