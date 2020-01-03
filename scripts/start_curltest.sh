#!/bin/bash

print_help() {
    echo ""
    echo "Usage:"
    MYNAME=`basename $0`
    echo "$MYNAME <goloba_port> <number_of_gets>"
    echo ""
    echo "Example:"
    echo "\$ $MYNAME 8000 3"
    echo ""
    echo "-  call 'curl -X GET localhost:8000' 3 times in a loop with 1s sleep in between"
}

echo ""
echo "Script for checking connection with GoLoBa using HTTP GET"
echo ""

if [ "$#" -ne 2 ];
then
    print_help
    exit 1
fi

# Check if curl installed
CURL_PATH=`which curl`
if [ "$?" != "0" ];
then
    echo "--"
    echo "Application 'curl' not found in PATH- exiting"
    echo "--"
    exit 1
fi

G_PORT=$1
N_TRIALS=$2

for i in `seq 1 1 $N_TRIALS`;
do
    echo ""
    echo "** Trial $i **"
    $CURL_PATH -X GET localhost:$G_PORT
    sleep 1
done

exit 0
