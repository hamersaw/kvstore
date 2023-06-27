#!/bin/bash

usage="USAGE: $(basename $0) [OPTIONS...]
OPTIONS:
    -h                 display this help message
    -o <op_count>      number of operations [default: 500]
    -s <start_id>      first id [default: 0]"

startid=0
ops=500
while getopts 'ho:s:' OPTION; do
    case "$OPTION" in 
        h)
            echo "$usage"
            exit 0
            ;;
        o)
            ops="$OPTARG"
            ;;
        s) 
            startid="$OPTARG"
            ;;
        ?) 
            echo "$usage"
            exit 1
            ;;
    esac
done

for i in $(seq $startid $(($startid + $ops - 1))); do
    curl -X POST -d "value=foo" localhost:3000/$i
done
