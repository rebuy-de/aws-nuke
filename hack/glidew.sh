#!/bin/bash

cd $( dirname $0 )/..
set -e

if [ ! -f target/glide/glide ]
then
	set -x
	mkdir -p target/glide

	VERSION=0.10.2
	FILE=glide-${VERSION}-$(go env GOHOSTOS)-$(go env GOHOSTARCH).tar.gz
	BASE=https://github.com/Masterminds/glide/releases/download
	URL=${BASE}/${VERSION}/${FILE}

	wget -c \
		-O target/glide/${FILE} \
		${URL}
	tar \
		-xzf target/glide/${FILE} \
		--strip-components=1 \
		-C target/glide

	set +x
fi

target/glide/glide "$@"
