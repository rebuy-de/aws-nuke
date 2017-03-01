# Source: https://github.com/rebuy-de/golang-template
# Version: 1.0.0
# Dependencies:
# * Glide
# * gocov (https://github.com/axw/gocov)
# * gocov-html (https://github.com/matm/gocov-html)

NAME=$(notdir $(PACKAGE))

BUILD_VERSION=$(shell git describe --always --dirty | tr '-' '.' )
BUILD_DATE=$(shell date -Iseconds)
BUILD_HASH=$(shell git rev-parse HEAD)
BUILD_MACHINE=$(shell echo $$HOSTNAME)
BUILD_USER=$(shell whoami)

BUILD_FLAGS=-ldflags "\
	-X '$(PACKAGE)/cmd.BuildVersion=$(BUILD_VERSION)' \
	-X '$(PACKAGE)/cmd.BuildDate=$(BUILD_DATE)' \
	-X '$(PACKAGE)/cmd.BuildHash=$(BUILD_HASH)' \
	-X '$(PACKAGE)/cmd.BuildEnvironment=$(BUILD_USER)@$(BUILD_MACHINE)' \
"

GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOPKGS=$(shell glide nv)

default: build

glide.lock: glide.yaml
	glide update

vendor: glide.lock glide.yaml
	glide install
	touch vendor

format:
	gofmt -s -w $(GOFILES)

test: vendor
	go test $(GOPKGS)

vet:
	go vet $(GOPKGS)

cov:
	gocov test -v $(GOPKGS) \
		| gocov-html > coverage.html

build: vendor
	go build \
		$(BUILD_FLAGS) \
		-o $(NAME)-$(BUILD_VERSION)
	ln -sf $(NAME)-$(BUILD_VERSION) $(NAME)


install: test
	go install \
		$(BUILD_FLAGS)

clean:
	rm -f $(NAME)*

.PHONY: build install test
