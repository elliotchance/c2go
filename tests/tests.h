#include <string.h> // strlen()
#include <math.h> // signbit()
#include <stdio.h> // printf()

// When comparing floating-point numbers this is how many significant bits will
// be used to calculate the epsilon.
#define INT64 52 - 4
#define INT32 23 - 4

// approx() is used to compare floating-point numbers. It will automatically
// calculate a very small epsilon (the difference between the two values that is
// still considered equal) based on the expected value and number of bits.
//
// approx() is used by the is_eq() macro for comparing numbers of any kind
// which makes it ideal for comparing different sized floats or other math that
// might introduce rounding errors.
#define approx(actual, expected) approxf(actual, expected, \
    sizeof(actual) != sizeof(expected) ? INT32 : INT64)

// approxf() will strictly not handle any input value that is infinite or NaN.
// It will always return false, even if the values are exactly equal. This is to
// force you to use the correct matcher (ie. is_inf()) instead of relying on
// comparisons which might not work.
static int approxf(double actual, double expected, int bits) {
    // We do not handle infinities or NaN.
    if (isinf(actual) || isinf(expected) || isnan(actual) || isnan(expected)) {
        return 0;
    }

    // If we expect zero (a common case) we have a fixed epsilon from actual. If
    // allowed to continue the epsilon calculated would be zero and we would be
    // doing an exact match which is what we want to avoid.
    if (expected == 0.0) {
        return fabs(actual) < (1 / pow(2, bits));
    }

    // The epsilon is calculated based on significant bits of the actual value.
    // The amount of bits used depends on the original size of the float (in
    // terms of bits) subtracting a few to allow for very slight rounding
    // errors.
    double epsilon = fabs(expected / pow(2, bits));

    // The numbers are considered equal if the absolute difference between them
    // is less than the relative epsilon.
    return fabs(actual - expected) <= epsilon;
}

// isnegzero tests if a value is a negative zero. Negative zeros are a special
// value of floating points that are equivalent in value, but not of bits to the
// common 0.0 value.
static int isnegzero(double x) {
    return (x * -0.0) == 0.0 && signbit(x);
}

// TAP: Below is a crude implementation of the TAP protocol for testing. If you
// add any new functionality to this file that is not part of TAP you should
// place that code above this comment.

// The number for the current test. It is incremented before each check.
static int current_test = 0;

// The total number of tests expected to be run. This must be set with plan()
// before any checks are performed.
static int total_tests = 0;

// The total number of failed tests.
static int total_failures = 0;

// Signals if the last test passed (1) or failed (0).
static int last_test_was_ok = 1;

// Set the number of expected tests/checks. This is important to prevent errors
// that would cause the program to silently fail before the test suite is
// finished.
//
// If the number of checks at the end (done_testing) is not the same as the
// value provided here then the suite will fail.
//
// To avoid Gcc warning (unused-function),
// force to call here the function isnegzero()
#define plan(numberOfTests)      \
    isnegzero(1);                \
    total_tests = numberOfTests; \
    printf("1..%d\n", numberOfTests)

// Print a diagnotic message or comment. Comments make the output more
// understandable (such as headings or descriptions) and can be printed at any
// time.
//
// diag() takes the same format and arguments as a printf().
#define diag(...)        \
    printf("# ");        \
    printf(__VA_ARGS__); \
    printf("\n");

// Check if the input is true. True in C is any value that is not a derivative
// of zero (or NULL). Any other value, including negative values and fractions
// are considered true.
//
// You should only use this when testing values that are integers or pointers.
// Zero values in floating-points are not always reliable to compare exactly.
#define is_true(actual)     \
    if (actual)             \
    {                       \
        pass("%s", #actual) \
    }                       \
    else                    \
    {                       \
        fail("%s", #actual) \
    }

// Check if the input is false. This works in the exact opposite way as
// is_true(). Be careful with floating-point values.
#define is_false(actual)    \
    if (actual == 0)        \
    {                       \
        pass("%s", #actual) \
    }                       \
    else                    \
    {                       \
        fail("%s", #actual) \
    }

// Immediately pass a check. This is useful if you are testing code flow, such
// as reaching a particular if/else branch.
//
// pass() takes the same arguments as a printf().
#define pass(...)                         \
    {                                     \
        ++current_test;                   \
        printf("%d ok - ", current_test); \
        printf(__VA_ARGS__);              \
        printf("\n");                     \
        last_test_was_ok = 1;             \
    }

// Immediately fail a check. This is useful if you are testing code flow (such
// as code that should not be reached or known bad conditions we met.
//
// fail() takes the same arguments as a printf().
#define fail(...)                                                        \
    {                                                                    \
        ++current_test;                                                  \
        ++total_failures;                                                \
        printf("%d not ok - %s:%d: ", current_test, __FILE__, __LINE__); \
        printf(__VA_ARGS__);                                             \
        printf("\n");                                                    \
        last_test_was_ok = 0;                                            \
    }

// Test if two values are equal. is_eq() will accept any *numerical* value than
// can be cast to a double and compared. This does not work with strings, see
// is_streq().
//
// Floating-point values are notoriously hard to compare exactly so the values
// are compared through the approx() function also defined in this file.
#define is_eq(actual, expected)                                            \
    if (approx((actual), (expected)))                                      \
    {                                                                      \
        pass("%s == %s", #actual, #expected)                               \
    }                                                                      \
    else                                                                   \
    {                                                                      \
        fail("%s == %s # got %.25g", #actual, #expected, (double)(actual)) \
    }

// This works in the opposite way as is_eq().
#define is_not_eq(actual, expected)                                     \
    if (approx((actual), (expected)))                                   \
    {                                                                   \
        fail("%s != %s # got %f", #actual, #expected, (double)(actual)) \
    }                                                                   \
    else                                                                \
    {                                                                   \
        pass("%s != %s", #actual, #expected)                            \
    }

// Compare two C strings. The two strings are equal if they are the same length
// (based on strlen) and each of their characters (actually bytes) are exactly
// the same value. This is nor to be used with strings that are mixed case or
// contain multibyte characters (eg. UTF-16, etc).
#define is_streq(actual, expected)                                 \
    if (strcmp(actual, expected) == 0)                             \
    {                                                              \
        pass("%s == %s", #actual, #expected)                       \
    }                                                              \
    else                                                           \
    {                                                              \
        fail("%s (%d b) == %s (%d b) # got \"%s\"", #actual,       \
            strlen(#actual), #expected, strlen(#expected), actual) \
    }

// Check that a floating-point value is Not A Number. Passing a value that is
// not floating point may lead to undefined behaviour.
#define is_nan(actual)                              \
    if (isnan(actual))                              \
    {                                               \
        pass("isnan(%s)", #actual)                  \
    }                                               \
    else                                            \
    {                                               \
        fail("isnan(%s) # got %f", #actual, actual) \
    }

// Check that a floating-point value is positive or negative infinity. A sign of
// less than 0 will check for -inf and a sign greater than 0 will check for
// +inf. A sign of 0 will always cause a failure.
//
// Mac and Linux check for infinity in different ways; either by using isinf()
// or comparing directly with an INFINITY constant. We always use both methods.
// This is separate checking the sign of the infinity.
#define is_inf(actual, sign)                                                                    \
    {                                                                                           \
        int linuxInf = ((actual) == -INFINITY || (actual) == INFINITY);                         \
        int macInf = (isinf(actual) == 1);                                                      \
        if ((linuxInf || macInf) && ((sign > 0 && (actual) > 0) || (sign < 0 && (actual) < 0))) \
        {                                                                                       \
            pass("isinf(%s, %d)", #actual, sign)                                                \
        }                                                                                       \
        else                                                                                    \
        {                                                                                       \
            fail("isinf(%s, %d) # got %f", #actual, sign, actual)                               \
        }                                                                                       \
    }

// Check that a value is a negative-zero. See isnegzero() for more information.
// Using non-floating-point values with this may give unexpected results.
#define is_negzero(actual) is_true(isnegzero(actual));

#define is_null(actual)             \
    if (actual == NULL)             \
    {                               \
        pass("%s == NULL", #actual) \
    }                               \
    else                            \
    {                               \
        fail("%s == NULL", #actual) \
    }

#define is_not_null(actual)         \
    if (actual != NULL)             \
    {                               \
        pass("%s != NULL", #actual) \
    }                               \
    else                            \
    {                               \
        fail("%s != NULL", #actual) \
    }

// disabled will suppress checks. The check, and any code associated within it
// will not be run at all. This is not the same as passing or skipping a test.
#define disabled(x)

// To signal that testing is now complete and to return the appropriate status
// code. This should be the last line in the main() function.
#define done_testing()                                                             \
    int exit_status = 0;                                                           \
    if (total_failures > 0)                                                        \
    {                                                                              \
        diag("FAILED: There was %d failed tests.", total_failures);                \
        exit_status = 101;                                                         \
    }                                                                              \
    if (current_test != total_tests)                                               \
    {                                                                              \
        diag("FAILED: Expected %d tests, but ran %d.", total_tests, current_test); \
        exit_status = 102;                                                         \
    }                                                                              \
    /* If we exit (with any status) the Go code coverage will not be generated. */ \
    if (exit_status != 0) {                                                        \
        return exit_status;                                                        \
    }                                                                              \
    return 0;

// or_return will return (with an optional value provided) if the check failed.
//
//     is_not_null(filePointer) or_return();
//
#define or_return(...)      \
    if (!last_test_was_ok)  \
    {                       \
        return __VA_ARGS__; \
    }

// Do not place code beyond this. See the TAP comment above.
