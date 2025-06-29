# embed-confluence

Provides shared C library that controls the [anacrolix/confluence](https://github.com/anacrolix/confluence) server for providing torrent contents over http.

## Example in C

- [main.c](_examples/c/main.c)
- [CMakeLists.txt](_examples/c/CMakeLists.txt)

## API

### Types

#### mbed_confluence_start_result

```c
typedef struct mbed_confluence_start_result {
    int32_t status;
    char* listen_addr;
    char* error;
} mbed_confluence_start_result;
```

#### mbed_confluence_result

```c
typedef struct mbed_confluence_result {
    int32_t status;
    char*  error;
} mbed_confluence_result;
```

### Methods

#### mbed_confluence_server_start

```c
mbed_confluence_start_result mbed_confluence_server_start(char* config, GoInt32 timeoutSec);
```

- `config` is a JSON config of the server, example:

    ```json
    {
        "addr": "127.0.0.1:8080",
        "public_ip_4": "40.114.177.156",
        "public_ip_6": null,
        "implicit_trackers": null,
        "torrent_grace_seconds": 60,
        "cache_capacity_bytes": 1073741824,
        "cache_dir":"cachefiles"
    }
    ```

#### mbed_confluence_server_wait

```c
void mbed_confluence_server_wait();
```

#### mbed_confluence_server_stop

```c
mbed_confluence_result mbed_confluence_server_stop();
```

#### free_mbed_confluence_start_result

```c
void free_mbed_confluence_start_result(mbed_confluence_start_result m);
```

#### free_mbed_confluence_result

```c
void free_mbed_confluence_result(mbed_confluence_result m);
```
