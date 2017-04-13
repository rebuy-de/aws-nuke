# Source: https://github.com/rebuy-de/golang-template
# Version: 1.3.1
# Dependencies:
# * Glide
# * gocov (https://github.com/axw/gocov)
# * gocov-html (https://github.com/matm/gocov-html)

NAME=$(notdir $(PACKAGE))

BUILD_VERSION=$(shell git describe --always --dirty --tags | tr '-' '.' )
BUILD_DATE=$(shell date)
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

vet:
	go vet $(GOPKGS)

lint:
	$(foreach pkg,$(GOPKGS),golint $(pkg);)

test_gopath:
	test $$(go list) = "$(PACKAGE)"

test_packages: vendor
	go test $(GOPKGS)

test_format:
	gofmt -l $(GOFILES)

test: test_gopath test_format vet lint test_packages

cov:
	gocov test -v $(GOPKGS) \
		| gocov-html > coverage.html

build: vendor
	go build \
		$(BUILD_FLAGS) \
		-o $(NAME)-$(BUILD_VERSION)-$(shell go env GOOS)-$(shell go env GOARCH)
	ln -sf $(NAME)-$(BUILD_VERSION)-$(shell go env GOOS)-$(shell go env GOARCH) $(NAME)

xc:
	GOOS=linux GOARCH=amd64 make build
	GOOS=darwin GOARCH=amd64 make build

install: test
	go install \
		$(BUILD_FLAGS)

clean:
	rm -f $(NAME)*

.PHONY: build install test
