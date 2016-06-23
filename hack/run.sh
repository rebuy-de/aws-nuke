#/bin/bash

source $( dirname $0)/deps.sh

go run *.go "$@"
