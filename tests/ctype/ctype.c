/* isalnum example */
#include <stdio.h>
#include <ctype.h>

void test_isalnum()
{
  int i;
  char str[] = "c3po...";
  i = 0;
  while (isalnum(str[i]))
  {
    i++;
  }
  printf("The first %d characters are alphanumeric.\n", i);
}

void test_isalpha()
{
  int i = 0;
  char str[] = "C++";
  while (str[i])
  {
    if (isalpha(str[i]))
      printf("character %c is alphabetic\n", str[i]);
    else
      printf("character %c is not alphabetic\n", str[i]);
    i++;
  }
}

void test_iscntrl()
{
  int i = 0;
  char str[] = "first line \n second line \n";
  while (!iscntrl(str[i]))
  {
    putchar(str[i]);
    i++;
  }
}

void test_isdigit()
{
  char str[] = "1776ad";
  int year;
  if (isdigit(str[0]))
  {
    year = atoi(str);
    printf("The year that followed %d was %d.\n", year, year + 1);
  }
}

void test_isgraph()
{
  int var1 = '3';
  int var2 = 'm';

  if (isgraph(var1))
  {
    printf("var1 = |%c| can be printed\n", var1);
  }
  else
  {
    printf("var1 = |%c| can't be printed\n", var1);
  }

  if (isgraph(var2))
  {
    printf("var2 = |%c| can be printed\n", var2);
  }
  else
  {
    printf("var2 = |%c| can't be printed\n", var2);
  }
}

void test_islower()
{
  int i = 0;
  char str[] = "Test String.\n";
  char c;
  while (str[i])
  {
    c = str[i];
    if (islower(c))
      c = toupper(c);
    putchar(c);
    i++;
  }
}

void test_isprint()
{
  int i = 0;
  char str[] = "first line \n second line \n";
  while (isprint(str[i]))
  {
    putchar(str[i]);
    i++;
  }
}

void test_ispunct()
{
  int i = 0;
  int cx = 0;
  char str[] = "Hello, welcome!";
  while (str[i])
  {
    if (ispunct(str[i]))
      cx++;
    i++;
  }
  printf("Sentence contains %d punctuation characters.\n", cx);
}

void test_isspace()
{
  char c;
  int i = 0;
  char str[] = "Example sentence to test isspace\n";
  while (str[i])
  {
    c = str[i];
    if (isspace(c))
      c = '\n';
    putchar(c);
    i++;
  }
}

void test_isupper()
{
  int i = 0;
  char str[] = "Test String.\n";
  char c;
  while (str[i])
  {
    c = str[i];
    if (isupper(c))
      c = tolower(c);
    putchar(c);
    i++;
  }
}

void test_isxdigit()
{
  char str[] = "ffff";
  long int number;
  if (isxdigit(str[0]))
  {
    number = strtol(str, NULL, 16);
    printf("The hexadecimal number %d is %d.\n", number, number);
  }
}

void test_tolower()
{
  int i = 0;
  char str[] = "Test String.\n";
  char c;
  while (str[i])
  {
    c = str[i];
    putchar(tolower(c));
    i++;
  }
}

void test_toupper()
{
  int i = 0;
  char str[] = "Test String.\n";
  char c;
  while (str[i])
  {
    c = str[i];
    putchar(toupper(c));
    i++;
  }
}

int main()
{
  test_isalnum();
  test_isalpha();
  test_iscntrl();
  test_isdigit();
  test_isgraph();
  test_islower();
  test_isprint();
  test_ispunct();
  test_isspace();
  test_isupper();
  test_isxdigit();
  test_tolower();
  test_toupper();
  
  return 0;
}
