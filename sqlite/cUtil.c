#include "cUtil.h"

/**
 * Converts a string of 'unsigned char' into 'char' then returns it as an allocated buffer.
 * Result must be freed
 */
char* ucharStrToCharStr(const unsigned char *input) {
    if(input == NULL) {
        fprintf(stderr, "C error: NULL passed to ucharStrToCharStr()\n");
        exit(1);
    }

    size_t len = strlen(input);
    char *res = malloc(len + 1);

    if(res == NULL) {
        fprintf(stderr, "C error: Failed to allocate result buffer in ucharStrToCharStr()");
        exit(1);
    }

    sprintf(res, "%s", input);
    
    return res;
}
