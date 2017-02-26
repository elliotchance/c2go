#include <stdio.h>
 
int main() 
{
    int array[100], n, c;
 
    printf("Enter the number of elements in array\n");
    scanf("%d", &n);
 
    printf("Enter %d elements\n", n);
 
    for ( c = 0 ; c < n ; c++ ) 
        scanf("%d", &array[c]);
 
    printf("Array elements entered by you are:\n");
 
    for ( c = 0 ; c < n ; c++ ) 
        printf("array[%d] = %d\n", c, array[c]);
 
    return 0;
}
