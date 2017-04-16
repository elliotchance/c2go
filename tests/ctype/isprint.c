/* isprint example */
#include <stdio.h>
#include <ctype.h>
int main ()
{
  int i=0;
  char str[]="first line \n second line \n";
  while (isprint(str[i]))
  {
    putchar (str[i]);
    i++;
  }
  return 0;
}
