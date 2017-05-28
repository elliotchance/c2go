#define RUN(t)                \
    printf("\n--- %s\n", #t); \
    t();

int approx(double actual, double expected)
{
    // The epsilon is calculated as one 5 millionths of the actual value. This
    // should be accurate enough, but also floating-points are usually rendered
    // with 6 places.
    double epsilon = fabs(expected * 0.00005);

    // The below line should be:
    //
    //     return fabs(a - b) <= ((fabs(a) < fabs(b) ? fabs(b) : fabs(a)) * epsilon);
    //
    // However, this is not yet supported. See:
    // https://github.com/elliotchance/c2go/issues/129
    double c = fabs(actual);
    if (fabs(actual) < fabs(expected))
    {
        c = fabs(expected);
    }

    return fabs(actual - expected) <= (c * epsilon);
}

// Test if x == -0.0.
int isnegzero(double x)
{
    return (x * -0.0) == 0.0 && signbit(x);
}

// The number for the current test.
static int current_test = 0;

// The total number of tests expected to be run.
static int total_tests = 0;

// The total number of failed tests.
static int total_failures = 0;

#define plan(numberOfTests)      \
    total_tests = numberOfTests; \
    printf("1..%d\n", numberOfTests)

#define diag(...)        \
    printf("# ");        \
    printf(__VA_ARGS__); \
    printf("\n");

#define is_true(actual)     \
    if (actual)             \
    {                       \
        pass("%s", #actual) \
    }                       \
    else                    \
    {                       \
        fail("%s", #actual) \
    }

#define pass(fmt, ...) \
    ++current_test;    \
    printf("%d ok - " fmt "\n", current_test, __VA_ARGS__);

#define fail(fmt, ...) \
    ++current_test;    \
    ++total_failures;  \
    printf("%d not ok - " fmt "\n", current_test, __VA_ARGS__);

#define is_eq(actual, expected)                              \
    if (approx((actual), (expected)))                        \
    {                                                        \
        pass("%s == %s", #actual, #expected)                 \
    }                                                        \
    else                                                     \
    {                                                        \
        fail("%s != %s, got %f", #actual, #expected, actual) \
    }

#define is_nan(actual)                                 \
    if (isnan(actual))                                 \
    {                                                  \
        pass("isnan(%s)", #actual)                     \
    }                                                  \
    else                                               \
    {                                                  \
        fail("%s is not NAN, got %f", #actual, actual) \
    }

#define is_inf(actual, sign)                                                              \
    if (isinf(actual) == 1 && ((sign > 0 && (actual) > 0) || (sign < 0 && (actual) < 0))) \
    {                                                                                     \
        pass("isinf(%s, %d)", #actual, sign)                                              \
    }                                                                                     \
    else                                                                                  \
    {                                                                                     \
        fail("%s is not +/-inf, got %d", #actual, isinf(actual))                          \
    }

#define is_negzero(actual) is_true(isnegzero(actual));

#define done_testing()                                                     \
    if (total_failures > 0)                                                \
    {                                                                      \
        diag("There was %d failed tests.", total_failures);                \
        return 101;                                                        \
    }                                                                      \
    if (current_test != total_tests)                                       \
    {                                                                      \
        diag("Expected %d tests, but ran %d.", total_tests, current_test); \
        return 102;                                                        \
    }                                                                      \
    return 0;
