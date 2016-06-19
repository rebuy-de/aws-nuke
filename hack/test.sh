#!/bin/bash

source $( dirname $0 )/deps.sh

go test -p 1 $(hack/glidew.sh novendor)
