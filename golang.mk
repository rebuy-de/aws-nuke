NAME=$(notdir $(PACKAGE))

BUILD_VERSION=$(shell git describe --always --dirty --tags | tr '-' '.' )
BUILD_DATE=$(shell date)
BUILD_HASH=$(shell git rev-parse HEAD)
BUILD_MACHINE=$(shell echo $$HOSTNAME)
BUILD_USER=$(shell whoami)

BUILD_FLAGS=-ldflags "\
	-s -w \
	-X '$(PACKAGE)/cmd.BuildVersion=$(BUILD_VERSION)' \
	-X '$(PACKAGE)/cmd.BuildDate=$(BUILD_DATE)' \
	-X '$(PACKAGE)/cmd.BuildHash=$(BUILD_HASH)' \
	-X '$(PACKAGE)/cmd.BuildEnvironment=$(BUILD_USER)@$(BUILD_MACHINE)' \
"

BUILD_ARTIFACT=$(NAME)-$(BUILD_VERSION)-$(shell go env GOOS)-$(shell go env GOARCH)

GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./.git/*")
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
		-o $(BUILD_ARTIFACT)
	ln -sf $(BUILD_ARTIFACT) $(NAME)

compress: build
	tar czf $(BUILD_ARTIFACT).tar.gz $(BUILD_ARTIFACT)

xc:
	GOOS=linux GOARCH=amd64 make compress
	GOOS=darwin GOARCH=amd64 make compress

install: vendor test
	go install \
		$(BUILD_FLAGS)

clean:
	rm -f $(NAME)*

.PHONY: build install test
