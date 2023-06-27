# kvstore
## Overview
TODO - channel vs map w/ lock


## Decisions
channel vs map w/ lock or syncmap: channel is more idiomatic golang - removes potential for deadlocks / liveness issues

channel implementation:
- everything in single request - alternatives are:
    - channel per request type - can not ensure requests are executed as they come
    - use internal interface with `isRequest` function so different request types work (similar to protobuf generation) - current use-case is simple, if more complicated should switch to something more verbose

- separate functions for "set" and "update" - arbitrary

- ideas: if slow we can use multiple shards / partitions where key is hashed and separate KVStore handle each shard - in best case, this allows # of concurrent operations equal to number of shards, can use `AtomicInt` to synchronize number of values
        

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
