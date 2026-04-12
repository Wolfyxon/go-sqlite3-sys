#pragma once

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/**
 * Converts a string of 'unsigned char' into 'char' then returns it as an allocated buffer.
 * Result must be freed
 */
char* ucharStrToCharStr(const unsigned char *input);