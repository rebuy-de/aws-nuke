FROM 074509403805.dkr.ecr.eu-west-1.amazonaws.com/rebuy-base-image-golang:latest
MAINTAINER reBuy reCommerce GmbH <dl-scb-tech-platform@rebuy.com>

COPY ["./", "/go/src/github.com/rebuy-de/aws-nuke/"]

RUN   set -ex && \
        cd /go/src/github.com/rebuy-de/aws-nuke && \
      glide install && \
      go test -v $(glide nv) && \
      go build -o /usr/bin/aws-nuke  github.com/rebuy-de/aws-nuke

ENTRYPOINT ["/usr/bin/aws-nuke"]
