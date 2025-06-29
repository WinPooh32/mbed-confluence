#ifndef MBED_CONFLUENCE_H
#define MBED_CONFLUENCE_H

#include <string.h>
#include <stdint.h>

typedef struct mbed_confluence_result {
    int32_t status;
    char*  error;
} mbed_confluence_result;

void free_mbed_confluence_result(mbed_confluence_result m);

typedef struct mbed_confluence_start_result {
    int32_t status;
    char* listen_addr;
    char* error;
} mbed_confluence_start_result;

void free_mbed_confluence_start_result(mbed_confluence_start_result m);

#endif