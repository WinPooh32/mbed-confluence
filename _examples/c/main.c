#include <stdio.h>
#include <signal.h> 
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>

#include "libmbedconfluence.h"

volatile sig_atomic_t exit_sig = 0;

void* stop_torrent_server_on_exit(void* arg) {
    while(!exit_sig){
        sleep(1);
    }

    mbed_confluence_result stop_res = mbed_confluence_server_stop();
    if (stop_res.status != 0) {
        printf("%s\n", stop_res.error);
        return NULL;
    }
    free_mbed_confluence_result(stop_res);

    printf("torrent server stopped!\n");
    
    return NULL;
}

void handle_signal(int sig)  { 
    printf("Caught signal %d\n", sig);

    if (sig == SIGINT) {
        exit_sig = sig;
    }
}

int main() {
    printf("start\n");

    pthread_t sig_stop_thread;
    pthread_create(&sig_stop_thread, NULL, stop_torrent_server_on_exit, NULL);
    signal(SIGINT, handle_signal); 

    char* cfg = R""""(
{
    "addr": "127.0.0.1:0",
    "public_ip_4": null,
    "public_ip_6": null,
    "implicit_trackers": null,
    "torrent_grace_seconds": 60,
    "cache_capacity_bytes": 1073741824,
    "cache_dir":"cachefiles"
}
)"""";

    mbed_confluence_start_result start_res = mbed_confluence_server_start(cfg, 60);
    if (start_res.status != 0) {
        printf("%s\n", start_res.error);
        return 1;
    }
    printf("server is listening at http://%s\n", start_res.listen_addr);
    free_mbed_confluence_start_result(start_res);

    mbed_confluence_server_wait();

    printf("exit\n");
    return exit_sig;
}