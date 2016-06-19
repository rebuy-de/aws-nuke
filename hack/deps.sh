#!/bin/bash

cd $( dirname $0 )/..
set -ex

export PROJECT=aws-nuke
export VERSION=$( git describe --always --dirty | tr '-' '.' )
export BINARY_NAME=${PROJECT}-${VERSION}

export PATH="$(readlink -f target/consul):${PATH}"

hack/glidew.sh install
