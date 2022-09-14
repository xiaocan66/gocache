#!/bin/bash

trap "rm server;kill 0" exit

go build -o server
./server -port=8001 &
./server -port=8002 &
./server -port=8003 &
./server -port=8004 &
./server -port=8005 -api=1 &
sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=tom" &
curl "http://localhost:9999/api?key=tom" &
curl "http://localhost:9999/api?key=tom" &
wait
