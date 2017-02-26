#import <ctype.h>
#import <stdio.h>

int main() {
    printf("var _DefaultRuneLocale _RuneLocale = _RuneLocale{\n\t__runetype: [256]uint32{");
    for (int i = 0; i < 256; ++i) {
        printf("%d, ", i, _DefaultRuneLocale.__runetype[i]);
    }
    printf("}\n}\n");
}
