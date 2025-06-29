#include <stdlib.h>
#include "mbed_confluence.h"

void free_mbed_confluence_result(mbed_confluence_result m){
    if (m.error != NULL) {
        free(m.error);
    }
}

void free_mbed_confluence_start_result(mbed_confluence_start_result m){
    if (m.listen_addr != NULL) {
        free(m.listen_addr);
    }
    if (m.error != NULL) {
        free(m.error);
    }
}