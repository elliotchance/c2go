// This program actually still works without including stdio.h but it should be
// here for consistency.

#include <stdio.h>
#include <string.h>
#include <assert.h>
#include "tests.h"

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_putchar()
{
    putchar('#');
    char c;
    for (c = 'A'; c <= 'Z'; c++)
        putchar(c);
    putchar('\n');

    pass("%s", "putchar");
}

void test_puts()
{
    puts("#c2go");

    pass("%s", "puts");
}

void test_printf()
{
    // TODO: printf() has a different syntax to Go
    // https://github.com/elliotchance/c2go/issues/94

    printf("# Characters: %c %c \n", 'a', 65);
    //printf("Decimals: %d %ld\n", 1977, 650000L);
    printf("# Preceding with blanks: %10d \n", 1977);
    printf("# Preceding with zeros: %010d \n", 1977);
    printf("# Some different radices: %d %x %o %#x %#o \n", 100, 100, 100, 100, 100);
    printf("# floats: %4.2f %+.0e %E \n", 3.1416, 3.1416, 3.1416);
    printf("# Width trick: %*d \n", 5, 10);
    printf("# %s \n", "A string");

    pass("%s", "printf");
}

void test_remove()
{
    // TODO: This does not actually test successfully deleting a file.
    if (remove("myfile.txt") != 0)
    {
        pass("%s", "error deleting file");
    }
    else
    {
        fail("%s", "file successfully deleted");
    }
}

void test_rename()
{
    // TODO: This does not actually test successfully renaming a file.
    int result;
    char oldname[] = "oldname.txt";
    char newname[] = "newname.txt";
    result = rename(oldname, newname);
    if (result == 0)
    {
        fail("%s", "File successfully renamed");
    }
    else
    {
        pass("%s", "Error renaming file");
    }
}

void test_fopen()
{
    FILE *pFile;
    pFile = fopen("/tmp/myfile.txt", "w");
    if (pFile != NULL)
    {
        is_not_null(pFile);
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
    is_not_null(pointer);
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
    pFile = fopen("/tmp/example.txt", "w+");
    is_not_null(pFile) or_return();

    fputs("test", pFile);
    fflush(pFile); // flushing or repositioning required
    fgets(mybuffer, 80, pFile);
    fclose(pFile);
}

void test_fprintf()
{
    FILE *pFile;
    int n;
    char *name = "John Smith";

    pFile = fopen("/tmp/myfile1.txt", "w");
    is_not_null(pFile);

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

    is_eq(f, 3.1416);
    is_streq(str, "PI");
}

void test_fgetc()
{
    FILE *pFile;
    int c;
    int n = 0;
    pFile = fopen("tests/stdio.c", "r");
    is_not_null(pFile);

    do
    {
        c = fgetc(pFile);
        if (c == '$')
            n++;
    } while (c != EOF);
    fclose(pFile);

    is_eq(n, 2);
}

void test_fgets()
{
    FILE *pFile;
    char *mystring;
    char dummy[20];

    pFile = fopen("tests/stdio.c", "r");
    is_not_null(pFile);

    mystring = fgets(dummy, 20, pFile);
    is_not_null(mystring);

    fclose(pFile);
}

void test_fputc()
{
    char c;

    fputc('#', stdout);
    for (c = 'A'; c <= 'Z'; c++)
        fputc(c, stdout);
    fputc('\n', stdout);
}

void test_fputs()
{
    FILE *pFile;
    char *sentence = "Hello, World";

    pFile = fopen("/tmp/mylog.txt", "w");
    fputs(sentence, pFile);
    fclose(pFile);
}

void test_getc()
{
    FILE *pFile;
    int c;
    int n = 0;
    pFile = fopen("tests/stdio.c", "r");
    is_not_null(pFile);

    do
    {
        c = getc(pFile);
        if (c == '$')
            n++;
    } while (c != EOF);
    fclose(pFile);

    is_eq(n, 2);
}

void test_putc()
{
    FILE *pFile;
    char c;

    pFile = fopen("/tmp/whatever.txt", "w");
    for (c = 'A'; c <= 'Z'; c++)
    {
        putc(c, pFile);
    }
    fclose(pFile);
}

void test_fseek()
{
    FILE *pFile;
    pFile = fopen("/tmp/example.txt", "w");
    fputs("This is an apple.", pFile);
    fseek(pFile, 9, SEEK_SET);
    fputs(" sam", pFile);
    fclose(pFile);
}

void test_ftell()
{
    FILE *pFile;
    long size;

    pFile = fopen("tests/stdio.c", "r");
    is_not_null(pFile);

    fseek(pFile, 0, SEEK_END); // non-portable
    size = ftell(pFile);
    fclose(pFile);

    is_eq(size, 7974);
}

void test_fread()
{
    FILE *pFile;
    int lSize;
    char buffer[1024];
    int result;

    pFile = fopen("tests/getchar.c", "r");
    is_not_null(pFile);

    // obtain file size:
    fseek(pFile, 0, SEEK_END);
    lSize = ftell(pFile);
    is_eq(lSize, 301);

    rewind(pFile);

    // copy the file into the buffer:
    result = fread(buffer, 1, lSize, pFile);
    is_eq(result, lSize);

    // See issue #107
    buffer[lSize - 1] = 0;

    is_eq(strlen(buffer), 300);

    // terminate
    fclose(pFile);
}

void test_fwrite()
{
    FILE *pFile;
    pFile = fopen("/tmp/myfile.bin", "w");
    fwrite("xyz", 1, 3, pFile);
    fclose(pFile);
}

void test_fgetpos()
{
    FILE *pFile;
    int c;
    int n;
    fpos_t pos;

    pFile = fopen("tests/stdio.c", "r");
    is_not_null(pFile);

    c = fgetc(pFile);
    is_eq(c, '/');

    fgetpos(pFile, &pos);
    for (n = 0; n < 3; n++)
    {
        fsetpos(pFile, &pos);
        c = fgetc(pFile);
        is_eq(c, '/');
    }

    fclose(pFile);
}

void test_fsetpos()
{
    FILE *pFile;
    fpos_t position;

    pFile = fopen("/tmp/myfile.txt", "w");
    fgetpos(pFile, &position);
    fputs("That is a sample", pFile);
    fsetpos(pFile, &position);
    fputs("This", pFile);
    fclose(pFile);
}

void test_rewind()
{
    int n;
    FILE *pFile;
    char buffer[27];

    pFile = fopen("/tmp/myfile.txt", "w+");
    for (n = 'A'; n <= 'Z'; n++)
        fputc(n, pFile);
    rewind(pFile);
    fread(buffer, 1, 26, pFile);
    fclose(pFile);
    buffer[26] = '\0';

    is_eq(strlen(buffer), 26);
}

void test_feof()
{
    FILE *pFile;
    int n = 0;
    pFile = fopen("tests/stdio.c", "r");
    is_not_null(pFile);

    while (fgetc(pFile) != EOF)
    {
        ++n;
    }
    if (feof(pFile))
    {
        pass("%s", "End-of-File reached.");
        is_eq(n, 7977);
    }
    else
    {
        fail("%s", "End-of-File was not reached.");
    }

    fclose(pFile);
}

int main()
{
    plan(33);

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
    START_TEST(fgetc)
    START_TEST(fgets)
    START_TEST(fputc)
    START_TEST(fputs)
    START_TEST(getc)
    START_TEST(putc)
    START_TEST(fseek)
    START_TEST(ftell)
    START_TEST(fread)
    START_TEST(fwrite)
    START_TEST(fgetpos)
    START_TEST(fsetpos)
    START_TEST(rewind)
    START_TEST(feof)

    done_testing();
}
