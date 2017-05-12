// This program actually still works without including stdio.h but it should be
// here for consistency.

#include <stdio.h>
#include <string.h>
#include <assert.h>

#define START_TEST(t)       \
    printf("--- %s\n", #t); \
    test_##t();

void test_putchar()
{
    char c;
    for (c = 'A'; c <= 'Z'; c++)
        putchar(c);
}

void test_puts()
{
    puts("c2go");
}

void test_printf()
{
    // TODO: printf() has a different syntax to Go
    // https://github.com/elliotchance/c2go/issues/94

    printf("Characters: %c %c \n", 'a', 65);
    //printf("Decimals: %d %ld\n", 1977, 650000L);
    printf("Preceding with blanks: %10d \n", 1977);
    printf("Preceding with zeros: %010d \n", 1977);
    printf("Some different radices: %d %x %o %#x %#o \n", 100, 100, 100, 100, 100);
    printf("floats: %4.2f %+.0e %E \n", 3.1416, 3.1416, 3.1416);
    printf("Width trick: %*d \n", 5, 10);
    printf("%s \n", "A string");
}

void test_remove()
{
    // TODO: This does not actually test successfully deleting a file.
    if (remove("myfile.txt") != 0)
        puts("Error deleting file");
    else
        puts("File successfully deleted");
}

void test_rename()
{
    // TODO: This does not actually test successfully renaming a file.
    int result;
    char oldname[] = "oldname.txt";
    char newname[] = "newname.txt";
    result = rename(oldname, newname);
    if (result == 0)
        puts("File successfully renamed");
    else
        puts("Error renaming file");
}

void test_fopen()
{
    FILE *pFile;
    pFile = fopen("/tmp/myfile.txt", "w");
    if (pFile != NULL)
    {
        fputs("fopen example", pFile);
        fclose(pFile);
    }
}

void test_tmpfile()
{
    char buffer[256];
    FILE *pFile;
    pFile = tmpfile();

    fputs("hello world", pFile);
    rewind(pFile);
    fputs(fgets(buffer, 256, pFile), stdout);
    fclose(pFile);
}

void test_tmpnam()
{
    // TODO: This is a tricky one to test because the output of tmpnam() in C
    // and Go will be different. I will keep the test here so at least it tries
    // to run the code but the test itself is not actually proving anything.

    char *pointer;

    // FIXME: We cannot pass variables by reference yet, which is a legitimate
    // way to use tmpnam(). I have to leave this disabled for now.
    //
    //     char buffer[L_tmpnam];
    //     tmpnam(buffer);
    //     assert(buffer != NULL);

    pointer = tmpnam(NULL);
    assert(pointer != NULL);
}

void test_fclose()
{
    FILE *pFile;
    pFile = fopen("/tmp/myfile.txt", "w");
    fputs("fclose example", pFile);
    fclose(pFile);
}

void test_fflush()
{
    char mybuffer[80];
    FILE *pFile;
    pFile = fopen("example.txt", "r+");
    if (pFile == NULL)
        printf("Error opening file");
    else
    {
        fputs("test", pFile);
        fflush(pFile); // flushing or repositioning required
        fgets(mybuffer, 80, pFile);
        puts(mybuffer);
        fclose(pFile);
    }
}

void test_fprintf()
{
    FILE *pFile;
    int n;
    char *name = "John Smith";

    pFile = fopen("/tmp/myfile1.txt", "w");
    assert(pFile != NULL);

    for (n = 0; n < 3; n++)
    {
        fprintf(pFile, "Name %d [%-10.10s]\n", n + 1, name);
    }

    fclose(pFile);
}

void test_fscanf()
{
    char str[80];
    float f;
    FILE *pFile;

    pFile = fopen("/tmp/myfile2.txt", "w+");
    fprintf(pFile, "%f %s", 3.1416, "PI");
    rewind(pFile);
    fscanf(pFile, "%f", &f);
    fscanf(pFile, "%s", str);
    fclose(pFile);
    printf("I have read: %f and %s \n", f, str);
}

void test_scanf()
{
    int i;

    scanf("%d", &i);
    printf("You enetered: %d\n", i);
}

void test_fgetc()
{
    FILE *pFile;
    int c;
    int n = 0;
    pFile = fopen("tests/stdio.c", "r");
    if (pFile == NULL)
        printf("Error opening file");
    else
    {
        do
        {
            c = fgetc(pFile);
            if (c == '$')
                n++;
        } while (c != EOF);
        fclose(pFile);
        printf("The file contains %d dollar sign characters ($).\n", n);
    }
}

void test_fgets()
{
    FILE *pFile;
    char *mystring;
    char dummy[20];

    pFile = fopen("tests/stdio.c", "r");
    if (pFile == NULL)
        printf("Error opening file");
    else
    {
        mystring = fgets(dummy, 20, pFile);
        if (mystring != NULL)
            puts(mystring);
        fclose(pFile);
    }
}

void test_fputc()
{
    char c;

    for (c = 'A'; c <= 'Z'; c++)
        fputc(c, stdout);
}

void test_fputs()
{
    FILE *pFile;
    char *sentence = "Hello, World";

    pFile = fopen("/tmp/mylog.txt", "w");
    fputs(sentence, pFile);
    fclose(pFile);
}

int main()
{
    START_TEST(putchar)
    START_TEST(puts)
    START_TEST(printf)
    START_TEST(remove)
    START_TEST(rename)
    START_TEST(fopen)
    START_TEST(tmpfile)
    START_TEST(tmpnam)
    START_TEST(fclose)
    START_TEST(fflush)
    START_TEST(printf)
    START_TEST(fprintf)
    START_TEST(fscanf)
    START_TEST(scanf)
    START_TEST(fgetc)
    START_TEST(fgets)
    START_TEST(fputc)
    START_TEST(fputs)

    return 0;
}
