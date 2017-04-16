[![Build Status](https://travis-ci.org/elliotchance/c2go.svg?branch=master)](https://travis-ci.org/elliotchance/c2go)
[![codecov](https://codecov.io/gh/elliotchance/c2go/branch/master/graph/badge.svg)](https://codecov.io/gh/elliotchance/c2go)

A tool for converting C to Go.

The goals of this project are:

1. To create a generic tool that can convert C to Go.
2. To be cross platform (linux and mac) and work against as many clang versions
as possible (the clang AST API is not stable).
3. To be a repeatable and predictable tool (rather than doing most of the work
and you have to clean up the output to get it working.)
4. To deliver quick and small version increments.
5. The ultimate milestone is to be able to compile the
[SQLite3 source code](https://sqlite.org/download.html) and have it working
without modification. This will be the 1.0.0 release.

# Installation

```bash
go get -u github.com/elliotchance/c2go
```

# Usage

```bash
c2go myfile.c
```

The `c2go` program processes a single C file and outputs the translated code
in Go. Let's use an included example,
[prime.c](https://github.com/elliotchance/c2go/blob/master/tests/math/prime.c):

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
c2go prime.c > prime.go
go run prime.go
```

```
Enter a number
23
Prime number.
```

`prime.go` looks like:

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

# What Is Supported?

See the
[Project Progress](https://github.com/elliotchance/c2go/wiki/Project-Progress).

# How It Works

This is the process:

1. The C code is preprocessed with clang. This generates a larger file (`pp.c`),
but removes all the platform specific directives and macros.

2. `pp.c` is parsed with the clang AST and dumps it in a colourful text format
that
[looks like this](http://ehsanakhgari.org/wp-content/uploads/2015/12/Screen-Shot-2015-12-03-at-5.02.38-PM.png).
Apart from just parsing the C and dumping an AST, the AST contains all of the
resolved information that a compiler would need (such as data types). This means
that the code must compile successfully under clang for the AST to also be
usable.

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

The test suite is run with `go test`.

# Contributing

Contributing is done with pull requests. There is no help that is too small! :)

If you're looking for where to start I can suggest
[finding a simple C program](http://www.programmingsimplified.com/c-program-examples)
(like the other examples) that do not successfully translate into Go.

Or, if you don't want to do that you can submit it as an issue so that it can be
picked up by someone else.
