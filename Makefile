
PACKAGE=github.com/rebuy-de/aws-nuke
VERSION=$(shell git describe --always --dirty | tr '-' '.' )
BUILD_FLAGS=-ldflags "-X $(PACKAGE)/cmd/version.version=$(VERSION)"

vendor: glide.lock glide.yaml
	glide install

build: vendor
	go build \
		$(BUILD_FLAGS) \
		-o aws-nuke-$(VERSION)

test: build
	go test $(shell glide nv)


install: test
	go install \
		$(BUILD_FLAGS)

.PHONY: build install
