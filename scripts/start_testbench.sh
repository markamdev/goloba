#!/bin/bash

TB_CONF_FILE="testbench.conf"
TB_LOG_FILE="test_$(date +"%Y%m%d_%H%M").log"

print_help() {
    echo ""
    echo "Usage:"
    MYNAME=`basename $0`
    echo "$MYNAME <goloba_port> <base_test_port> <number_of_test_servers>"
    echo ""
    echo "Example:"
    echo "\$ $MYNAME 7000 8000 3"
    echo ""
    echo "- launches 3 dummyserver instances listening on ports 8000, 8001 and 8002"
    echo "- prepares config for GoLoBa with 7000 as listenign port"
}

echo ""
echo "** Script for starting simple testbench for GoLoBa **"
echo ""

if [ "$#" -ne 3 ]; then
    print_help
    exit 1
fi

# Save commandline params to variables
G_PORT=$1
B_PORT=$2
N_SERV=$3

# Check where script has been called from 
# and where are goloba and dummy server binaries
CUR_DIR="$PWD"
SUB_DIR=`basename $CUR_DIR`
WORK_DIR=""
if [ "$SUB_DIR" == "scripts" ];
then
    WORK_DIR="$CUR_DIR/../build"
elif [ "$SUB_DIR" == "build" ];
then
    WORK_DIR="./"
elif [ "$SUB_DIR" == "goloba" ];
then
    WORK_DIR="$CUR_DIR/build"
else
    echo "Where am I? Please call this script from:"
    echo "- main goloba repo directory"
    echo "- ./scripts subdir"
    echo "- ./build subdir"
    exit 1
fi

echo "[Working directory set to: $WORK_DIR]"
cd $WORK_DIR

echo ""
echo "** Preparing testbench with $N_SERV listeting servers **"

PORTS=`seq --separator " " $B_PORT 1 $(( B_PORT + N_SERV - 1))`
echo "- ports to be used: $PORTS"

# Prepare config file beginning
echo "{" > $TB_CONF_FILE
echo "    \"port\":$G_PORT," >> $TB_CONF_FILE
echo -n "    \"servers\": [" >> $TB_CONF_FILE

# Launch $N_SERV instances of servers and add necessary config entries
for port in $PORTS
do
    echo "<> Launching dummyserver listeting on port $port <>"
    # Add comma before next server definition
    if [ "$port" != "$B_PORT" ];
    then
        echo "," >> $TB_CONF_FILE
    fi
    echo -n " \"localhost:$port\"" >> $TB_CONF_FILE
    ./dummyserver -p $port -m "TestBench listener at $port" &
    sleep 1
done

echo " ]" >> $TB_CONF_FILE
echo "}" >> $TB_CONF_FILE

# start GoLoBa with given config
echo "** Launching GoLoBa load balancer (should be listening on port $G_PORT) **"
./goloba -f $TB_CONF_FILE -l $TB_LOG_FILE

# Kill all servers if GoLoBa finished
echo "** GoLoBa process stopped - killing all listeners **"
killall dummyserver

exit 0
