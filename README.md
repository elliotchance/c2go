[![Build Status](https://travis-ci.org/elliotchance/c2go.svg?branch=master)](https://travis-ci.org/elliotchance/c2go)

A tool for converting C to Go.

# Why?

I created this project as a proof of concept. It is written in python and uses
the clang AST to do all of the hard work. Python was an easy was to get up and
running quickly. The intention is to convert the Python to Go.

The goals of this project are:

1. To create a generic tool that can convert C to Go.
2. To be cross platform (linux and mac) and work against all versions of clang
3+.
2. To be a repeatable and predicatble tool (rather than doing most of the work
and you have to clean up the output to get it working.)
3. To be written in Go (eventually).
4. To deliver quick and small version increments.
5. The ultimate milestone is to be able to compile the
[SQLite3 source code](https://sqlite.org/download.html) and have it working
without modification. This will be the 1.0.0 release.

# Installation

At the moment the easiest way to install it is to simply clone the repo:

```bash
git clone https://github.com/elliotchance/c2go.git
```

# Usage

The `c2go.py` returns a list of files that need to be compiled with Go:

```bash
python c2go.py tests/prime.c
```

```
out.go
functions.go
functions-Linux.go
```

These files depend on the `#include` and other options when converting the C.
You can send these files directly into `go run`:

```bash
go run $(python c2go.py tests/prime.c)
```

```
Enter a number
23
Prime number.
```

# What Is Supported?

This table represents what is supported. If you see anything missing (there is a
lot missing!) please add it with a pull request.

| Function      | Supported?    | Notes                      |
| ------------- | ------------- | -------------------------- |
| **math.h**    | Partly        | All of the C99 functions.  |
| acos          | Yes           |                            |
| asin          | Yes           |                            |
| atan          | Yes           |                            |
| atan2         | Yes           |                            |
| ceil          | Yes           |                            |
| cos           | Yes           |                            |
| cosh          | Yes           |                            |
| exp           | Yes           |                            |
| fabs          | Yes           |                            |
| floor         | Yes           |                            |
| fmod          | Yes           |                            |
| ldexp         | Yes           |                            |
| log           | Yes           |                            |
| log10         | Yes           |                            |
| pow           | Yes           |                            |
| sin           | Yes           |                            |
| sinh          | Yes           |                            |
| sqrt          | Yes           |                            |
| tan           | Yes           |                            |
| tanh          | Yes           |                            |
| **stdio.h**   | Partly        |                            |
| printf        | Yes           |                            |
| scanf         | Yes           |                            |

# How It Works

This is the process:

1. The C code is preprocessed with clang. This generates a larger file (`pp.c`),
but removes all the platform specific directives and macros.

2. `pp.c` is parsed with the clang AST and dumps it in a colourful text format that
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

The test suite is run with
[run-tests.sh](https://github.com/elliotchance/c2go/blob/master/run-tests.sh).

# Contributing

As I said it is still very early days (sorry for all the hacky Python). And
eventually the build chain can be converted to pure Go since we don't need any
clang APIs.

Contributing is done with pull requests. There is no help that is too small! :)
If you're looking for where to start I can suggest
[finding a simple C program](http://www.programmingsimplified.com/c-program-examples)
(like the other examples) that does not successful translate into Go and fixing
up the Python so that it does.

Or, if you don't want to do that you can submit it as an issue so that it can be
picked up by someone else.
