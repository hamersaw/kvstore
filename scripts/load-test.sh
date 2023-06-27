#!/bin/bash

usage="USAGE: $(basename $0) [OPTIONS...]
OPTIONS:
    -c <client_count>  number of clients [default: 5]
    -h                 display this help message
    -o <op_count>      number of operations [default: 500]"

clients=5
ops=500
while getopts 'c:ho:' OPTION; do
    case "$OPTION" in 
        c) 
            clients="$OPTARG"
            ;;
        h)
            echo "$usage"
            exit 0
            ;;
        o)
            ops="$OPTARG"
            ;;
        ?) 
            echo "$usage"
            exit 1
            ;;
    esac
done

for i in $(seq 1 $clients); do
    ./scripts/benchmark.sh -o $ops -s $(( ($i - 1) * $ops )) &
    pids[${i}]=$!
done

for pid in ${pids[*]}; do
    wait $pid
done
