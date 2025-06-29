package main

/*
#include "mbed_confluence.h"
*/
import "C"
import (
	"context"
	"encoding/json"
	"time"
)

var srv Server

//export mbed_confluence_server_start
func mbed_confluence_server_start(config *C.char, timeoutSec int32) C.mbed_confluence_start_result {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	var cfg Config
	if err := json.Unmarshal([]byte(C.GoString(config)), &cfg); err != nil {
		return C.mbed_confluence_start_result{
			status: 1,
			error:  C.CString(err.Error()),
		}
	}

	listenAddr, err := srv.Start(ctx, cfg)
	if err != nil {
		return C.mbed_confluence_start_result{
			status: 1,
			error:  C.CString(err.Error()),
		}
	}

	return C.mbed_confluence_start_result{
		status:      0,
		listen_addr: C.CString(listenAddr),
		error:       nil,
	}
}

//export mbed_confluence_server_wait
func mbed_confluence_server_wait() {
	srv.Wait()
}

//export mbed_confluence_server_stop
func mbed_confluence_server_stop() C.mbed_confluence_result {
	if err := srv.Stop(context.Background()); err != nil {
		return C.mbed_confluence_result{
			status: 1,
			error:  C.CString(err.Error()),
		}
	}

	return C.mbed_confluence_result{
		status: 0,
		error:  nil,
	}
}
