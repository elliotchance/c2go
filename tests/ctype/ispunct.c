/* ispunct example */
#include <stdio.h>
#include <ctype.h>
int main ()
{
  int i=0;
  int cx=0;
  char str[]="Hello, welcome!";
  while (str[i])
  {
    if (ispunct(str[i])) cx++;
    i++;
  }
  printf ("Sentence contains %d punctuation characters.\n", cx);
  return 0;
}