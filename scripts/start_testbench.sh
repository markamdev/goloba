#!/bin/bash

BUILD_DIR=./build
GOLOBA_BIN=$BUILD_DIR/goloba
DUMMY_BIN=$BUILD_DIR/dummyserver
TB_CONF_FILE=$BUILD_DIR/testbench.conf
TB_LOG_FILE=$BUILD_DIR/"test_$(date +"%Y%m%d_%H%M").log"

echo "Testing scenario"
echo "- 4 dummyserver instances launched, listening on ports 8001, 8002, 8003 and 8004"
echo "- GoLoBa starts listening on port 8000, redirects traffic to ports above"

if [[ ! -e $GOLOBA_BIN || ! -e $DUMMY_BIN ]];
then
    echo "--"
    echo "Missing binaries - did you build GoLoBa and Dummyserver?"
    echo "--"
    exit 1
fi

if [[ ! -x $GOLOBA_BIN || ! -x $DUMMY_BIN ]];
then
    echo "--"
    echo "GoLoBa or Dummyserver binary not executable. Exiting"
    echo "--"
    exit 1
fi

CONF_DATA='{
    "port":8000,
    "servers": [ "localhost:8001", "localhost:8002","localhost:8003","localhost:8004"]
}'
echo $CONF_DATA > $TB_CONF_FILE

# start GoLoBa with given config
echo "Launching GoLoBa load balancer listening on port 8000"
$GOLOBA_BIN -f $TB_CONF_FILE -l $TB_LOG_FILE &
GLB_PID=$!
echo "... launched as process $GLB_PID"

# start listening dummy servers
echo "Launching Dummyserver instances listening on ports 8001..8004"
./build/dummyserver -p 8001 -m "Testbench server 1" &
./build/dummyserver -p 8002 -m "Testbench server 2" &
./build/dummyserver -p 8003 -m "Testbench server 3" &
./build/dummyserver -p 8004 -m "Testbench server 4" &

# launch CURL Testing
echo "Launching simple CURL-based testing"
./scripts/start_curltest.sh 8000 10
if [ $? -ne 0 ];
then
    echo "!! TESTING FAILED !!"
    exit 1
fi

# kill GoLoBa and then server
echo "Stopping GoLoBa and servers"
kill $GLB_PID
killall dummyserver

exit 0
