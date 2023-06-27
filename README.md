# kvstore
## Overview
This KVStore implementation contains two approaches to handle concurrent requests within a KVStore,
namely using a request channel (CSP) and RWMutex to serialize request handling. Both approaches are
exposed under a RESTful interface.

## Usage
    # compile
    make build

    # start kvstore
    ./kvstore --concurrency channel
    ./kvstore --concurrency rwmutex

    # test REST API
    curl -X POST -d "value=bar" localhost:3000/foo -v
    curl -X GET localhost:3000/foo -v
    curl -X PUT -d "value=baz" localhost:3000/foo -v
    curl -X DELETE localhost:3000/foo -v

    # perf test
    time ./scripts/load-test.sh -c 10 -o 1000
