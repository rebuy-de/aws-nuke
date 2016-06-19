#/bin/bash

source $( dirname $0)/test.sh

mkdir -p target

go build \
	-ldflags "-X main.version=${VERSION}" \
	-o target/${BINARY_NAME}
