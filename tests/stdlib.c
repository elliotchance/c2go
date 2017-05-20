#include <stdio.h>
#include <stdlib.h>

void test_malloc()
{
    int i, n;
    char *buffer;

    printf("How long do you want the string? ");
    scanf("%d", &i);

    buffer = (char *)malloc(i + 1);
    if (buffer == NULL)
        return;

    for (n = 0; n < i; n++)
        buffer[n] = i % 26 + 'a';
    buffer[i] = '\0';

    printf("Random string: %s\n", buffer);
    free(buffer);
}

int main()
{
    test_malloc();
    
    return 0;
}
